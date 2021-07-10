package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	server "github.com/bartmika/tstorage-server/internal"
	"github.com/bartmika/tstorage-server/utils"
)

var (
	port                     int
	dataPath                 string
	timestampPrecision       string
	partitionDurationInHours int
	writeTimeoutInSeconds    int
)

func init() {
	// The following are optional and will have defaults placed when missing.
	serveCmd.Flags().IntVarP(&port, "port", "p", 50051, "The port to run this server on")
	serveCmd.Flags().StringVarP(&dataPath, "dataPath", "d", "./tsdb", "The location to save the database files to.")
	serveCmd.Flags().StringVarP(&timestampPrecision, "timestampPrecision", "t", "s", "The precision of timestamps to be used by all operations. Options: ")
	serveCmd.Flags().IntVarP(&partitionDurationInHours, "partitionDurationInHours", "b", 1, "The timestamp range inside partitions.")
	serveCmd.Flags().IntVarP(&writeTimeoutInSeconds, "writeTimeoutInSeconds", "w", 30, "The timeout to wait when workers are busy (in seconds).")

	// Make this sub-command part of our application.
	rootCmd.AddCommand(serveCmd)
}

func doServe() {
	// Convert the user inputted integer value to be a `time.Duration` type.
	partitionDuration := time.Duration(partitionDurationInHours) * time.Hour
	writeTimeout := time.Duration(writeTimeoutInSeconds) * time.Second

	// Setup our server.
	server := server.New(port, dataPath, timestampPrecision, partitionDuration, writeTimeout)

	// DEVELOPERS CODE:
	// The following code will create an anonymous goroutine which will have a
	// blocking chan `sigs`. This blocking chan will only unblock when the
	// golang app receives a termination command; therfore the anyomous
	// goroutine will run and terminate our running application.
	//
	// Special Thanks:
	// (1) https://gobyexample.com/signals
	// (2) https://guzalexander.com/2017/05/31/gracefully-exit-server-in-go.html
	//
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs // Block execution until signal from terminal gets triggered here.
		server.StopMainRuntimeLoop()
	}()
	server.RunMainRuntimeLoop()
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the gRPC server",
	Long:  `Run the gRPC server to allow other services to access the storage application`,
	Run: func(cmd *cobra.Command, args []string) {
		// Defensive code. Make sure the user selected the correct `timestampPrecision`
		// choices before continuing execution of our command.
		okTimestampPrecision := []string{"ns", "us", "ms", "s"}
		if utils.Contains(okTimestampPrecision, timestampPrecision) == false {
			log.Fatal("Timestamp precision must be either one of the following: ns, us, ms, or s.")
		}

		// Execute our command with our validated inputs.
		doServe()
	},
}
