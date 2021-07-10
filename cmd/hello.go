package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials"

	pb "github.com/bartmika/tstorage-server/proto"
)

var (
	name string
)

func init() {
	// The following are required.
	helloCmd.Flags().StringVarP(&name, "name", "n", "Anonymous", "The name to send the server.")
	helloCmd.MarkFlagRequired("name")

	// The following are optional and will have defaults placed when missing.
	helloCmd.Flags().IntVarP(&port, "port", "p", 50051, "The port of our server.")
	rootCmd.AddCommand(helloCmd)
}

func doHello() {
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

	// Perform our gRPC request.
	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// Print out the gRPC response.
	log.Printf("Server Response: %s", r.GetMessage())
}

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Send hello message to gRPC server",
	Long:  `Connect to the gRPC server and send a hello message. Command used to test out that the server is running.`,
	Run: func(cmd *cobra.Command, args []string) {
		doHello()
	},
}
