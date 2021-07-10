package internal

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/nakabonne/tstorage"
	"google.golang.org/grpc"

	pb "github.com/bartmika/tstorage-server/proto"
)

type TStorageServer struct {
	port               int
	dataPath           string
	timestampPrecision tstorage.TimestampPrecision
	partitionDuration  time.Duration
	writeTimeout       time.Duration
	storage            tstorage.Storage
	grpcServer         *grpc.Server
}

func New(port int, dataPath string, timestampPrecision string, partitionDuration time.Duration, writeTimeout time.Duration) *TStorageServer {
	// Conver to the format that is accepted by the library.
	var tsp tstorage.TimestampPrecision
	switch timestampPrecision {
	case "ns":
		tsp = tstorage.Nanoseconds
	case "us":
		tsp = tstorage.Microseconds
	case "ms":
		tsp = tstorage.Milliseconds
	case "s":
		tsp = tstorage.Seconds
	}

	return &TStorageServer{
		port:               port,
		dataPath:           dataPath,
		timestampPrecision: tsp,
		partitionDuration:  partitionDuration,
		writeTimeout:       writeTimeout,
		storage:            nil,
		grpcServer:         nil,
	}
}

// Function will consume the main runtime loop and run the business logic
// of the application.
func (s *TStorageServer) RunMainRuntimeLoop() {
	// Open a TCP server to the specified localhost and environment variable
	// specified port number.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Initialize our gRPC server using our TCP server.
	grpcServer := grpc.NewServer()

	// Initialize our fast time-series database.
	storage, _ := tstorage.NewStorage(
		tstorage.WithDataPath(s.dataPath),
		tstorage.WithTimestampPrecision(s.timestampPrecision),
		tstorage.WithPartitionDuration(s.partitionDuration),
		tstorage.WithWriteTimeout(s.writeTimeout),
	)

	// Save reference to our application state.
	s.grpcServer = grpcServer
	s.storage = storage

	// For debugging purposes only.
	log.Printf("gRPC server is running.")

	// Block the main runtime loop for accepting and processing gRPC requests.
	pb.RegisterTStorageServer(grpcServer, &TStorageServerImpl{
		// DEVELOPERS NOTE:
		// We want to attach to every gRPC call the following variables...
		storage: s.storage,
	})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Function will tell the application to stop the main runtime loop when
// the process has been finished.
func (s *TStorageServer) StopMainRuntimeLoop() {
	log.Printf("Starting graceful shutdown now...")

	// Finish our database operations running.
	s.storage.Close()

	// Finish any RPC communication taking place at the moment before
	// shutting down the gRPC server.
	s.grpcServer.GracefulStop()
}
