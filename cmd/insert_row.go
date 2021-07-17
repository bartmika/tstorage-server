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

var (
	metric string
	value  float64
	tsv    int64
)

func init() {
	// The following are required.
	insertRowCmd.Flags().StringVarP(&metric, "metric", "m", "", "The metric to attach to the TSD.")
	insertRowCmd.MarkFlagRequired("metric")
	insertRowCmd.Flags().Float64VarP(&value, "value", "v", 0.00, "The value to attach to the TSD.")
	insertRowCmd.MarkFlagRequired("value")
	insertRowCmd.Flags().Int64VarP(&tsv, "timestamp", "t", 0, "The timestamp to attach to the TSD.")
	insertRowCmd.MarkFlagRequired("timestamp")

	// The following are optional and will have defaults placed when missing.
	insertRowCmd.Flags().IntVarP(&port, "port", "p", 50051, "The port of our server.")
	rootCmd.AddCommand(insertRowCmd)
}

func doInsertRow() {
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

	// Perform our gRPC request.
	_, err = client.InsertRow(ctx, &pb.TimeSeriesDatum{Labels: labels, Metric: metric, Value: value, Timestamp: ts})
	if err != nil {
		log.Fatalf("could not add: %v", err)
	}

	log.Printf("Successfully inserted")
}

var insertRowCmd = &cobra.Command{
	Use:   "insert_row",
	Short: "Insert single datum",
	Long:  `Connect to the gRPC server and sends a single time-series datum.`,
	Run: func(cmd *cobra.Command, args []string) {
		doInsertRow()
	},
}
