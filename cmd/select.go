package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials"

	tspb "github.com/golang/protobuf/ptypes/timestamp"

	pb "github.com/bartmika/tstorage-server/proto"
)

var (
	start int64
	end   int64
)

func init() {
	// The following are required.
	selectCmd.Flags().StringVarP(&metric, "metric", "m", "", "The metric to filter by")
	selectCmd.MarkFlagRequired("metric")
	selectCmd.Flags().Int64VarP(&start, "start", "s", 0, "The start timestamp to begin our range")
	selectCmd.MarkFlagRequired("start")
	selectCmd.Flags().Int64VarP(&end, "end", "e", 0, "The end timestamp to finish our range")
	selectCmd.MarkFlagRequired("end")

	// The following are optional and will have defaults placed when missing.
	selectCmd.Flags().IntVarP(&port, "port", "p", 50051, "The port of our server.")
	rootCmd.AddCommand(selectCmd)
}

func doSelectRow() {
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

	// Convert the unix timestamp into the protocal buffers timestamp format.
	sts := &tspb.Timestamp{
		Seconds: start,
		Nanos:   0,
	}
	ets := &tspb.Timestamp{
		Seconds: end,
		Nanos:   0,
	}

	// Generate our labels.
	labels := []*pb.Label{}
	labels = append(labels, &pb.Label{Name: "Source", Value: "Command"})

	// Perform our gRPC request.
	stream, err := client.Select(ctx, &pb.Filter{Labels: labels, Metric: metric, Start: sts, End: ets})
	if err != nil {
		log.Fatalf("could not select: %v", err)
	}

	// Handle our stream of data from the server.
	for {
		dataPoint, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error with stream: %v", err)
		}

		// Print out the gRPC response.
		log.Printf("Server Response: %s", dataPoint)
	}
}

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "List data",
	Long:  `Connect to the gRPC server and return list of results based on a selection filter.`,
	Run: func(cmd *cobra.Command, args []string) {
		doSelectRow()
	},
}
