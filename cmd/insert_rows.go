package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials"

	tspb "github.com/golang/protobuf/ptypes/timestamp"

	pb "github.com/bartmika/tstorage-server/proto"
)

func init() {
	// The following are required.
	insertRowsCmd.Flags().StringVarP(&metric, "metric", "m", "", "The metric to attach to the TSD.")
	insertRowsCmd.MarkFlagRequired("metric")
	insertRowsCmd.Flags().Float64VarP(&value, "value", "v", 0.00, "The value to attach to the TSD.")
	insertRowsCmd.MarkFlagRequired("value")
	insertRowsCmd.Flags().Int64VarP(&tsv, "timestamp", "t", 0, "The timestamp to attach to the TSD.")
	insertRowsCmd.MarkFlagRequired("timestamp")

	// The following are optional and will have defaults placed when missing.
	insertRowsCmd.Flags().IntVarP(&port, "port", "p", 50051, "The port of our server.")
	rootCmd.AddCommand(insertRowsCmd)
}

func doInsertRows() {
	// Set up a direct connection to the gRPC server.
	conn, err := grpc.Dial(
		fmt.Sprintf(":%v", port),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// Set up our protocol buffer interface.
	client := pb.NewTStorageClient(conn)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ts := &tspb.Timestamp{
		Seconds: tsv,
		Nanos:   0,
	}

	// Generate our labels.
	labels := []*pb.Label{}
	labels = append(labels, &pb.Label{Name: "Source", Value: "Command"})

	stream, err := client.InsertRows(ctx)
	if err != nil {
		log.Fatalf("%v.InsertRows(_) = _, %v", client, err)
	}

	tsd := &pb.TimeSeriesDatum{Labels: labels, Metric: metric, Value: value, Timestamp: ts}

	// DEVELOPERS NOTE:
	// To stream from a client to a server using gRPC, the following documentation
	// will help explain how it works. Please visit it if the code below does
	// not make any sense.
	// https://grpc.io/docs/languages/go/basics/#client-side-streaming-rpc-1

	if err := stream.Send(tsd); err != nil {
		log.Fatalf("%v.Send(%v) = %v", stream, tsd, err)
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("Server Response: %v", reply)
}

var insertRowsCmd = &cobra.Command{
	Use:   "insert_rows",
	Short: "Insert single datum using streaming",
	Long:  `Connect to the gRPC server and send a time-series datum using the streaming RPC.`,
	Run: func(cmd *cobra.Command, args []string) {
		doInsertRows()
	},
}
