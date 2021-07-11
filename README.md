# tstorage-server
Persistent fast time-series data storage database accessible over gRPC

`tstorage-server` is lightweight local on-disk storage engine server for time-series data accessible over gRPC. Run this server once and share fast time-series data CRUD operations in either local or remote applications.
The purpose of this server is to allow interprocess communication overtop the [`tstorage`](https://github.com/nakabonne/tstorage) package.

## Prerequisites

You must have the following installed before proceeding. If you are missing any one of these then you cannot begin.

* ``Go 1.16.3``

## Installation

Get our latest code.

```bash
go install github.com/bartmika/tstorage-server
```

## Usage

To start the server, run the following command in your **terminal**:

```bash
$GOBIN/tstorage-server serve
```

That's it! If everything works, you should see a message saying `gRPC server is running.`.

## Subcommands Reference

### ``serve``
Run the gRPC server to allow other services to access the storage application

Fields

* `-p` or `--port` is for the port for this server to run on. Default value is 50051 if you don't use this option.
* `-d` or `--dataPath` is for the location path to use to save the database files to. Default value is './tsdb'.
* `-t` or `--timestampPrecision` is used for the precision of timestamps on all operations. The available options are "ns", "us", "ms", "s". Default value is "s".
* `-b` or `--partitionDurationInHours` is used for the timestamp range inside partitions. Default value is 1.
* `-w` or `--writeTimeoutInSeconds` is for timeout to wait when workers are busy (in seconds). Default value is 30.

Example:

```bash
$GOBIN/tstorage-server serve -p=50051 -d="./tsdb" -t="s" -b=1 -w=30
```

### ``version``
Prints the current version of our server

Example:

```bash
$GOBIN/tstorage-server version
```

### ``insert_row``
Connect to the running gRPC server and sends a single time-series datum to insert. Please note, you need to have your `tstorage-server` running with the `serve` subcommand for this `insert_row` command to work!

Fields

* `-p` or `--port` is for the port for this server to run on. Default value is 50051 if you don't use this option.
* `-m` or `--metric` is for the metric to insert with the time-series datum. This field is required.
* `-v` or `--value` is for the value of the time-series datum. This field is required.
* `-t` or `--timestamp` is for the unix epoch time for the time-series datum. This field is required.

Example:

```bash
$GOBIN/tstorage-server insert_row -p=50051 -m="solar_biodigester_temperature_in_degrees" -v=50 -t=1600000000
```

Developer Notes:
- There also exists a `insert_rows` subcommand but it works exactly as `insert_row` command with the exception that the internal code is using streaming. This is done so programmers can look at the code and see how to use streaming of time-series data.

### ``select``
Connect to the gRPC server and return list of results based on a selection filter.

Fields

* `-p` or `--port` is for the port for this server to run on. Default value is 50051 if you don't use this option.
* `-m` or `--metric` is for the metric to insert with the time-series datum. This field is required.
* `-s` or `--start` is for the start timestamp to begin our range by.
* `-e` or `--end` is for the end timestamp to finish our range by.

Example:

```bash
$GOBIN/tstorage-server select --port=50051 --metric="bio_reactor_pressure_in_kpa" --start=1600000000 --end=1725946120
```


## How to Access using gRPC

### Example 1 - Insert a Single Row

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	tspb "github.com/golang/protobuf/ptypes/timestamp"

	pb "github.com/bartmika/tstorage-server/proto"
)

func main() {
  port := 50051

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
	labels = append(labels, &pb.Label{Name: "Host", Value: "127.0.0.1"})

	// Perform our gRPC request.
	r, err := client.InsertRow(ctx, &pb.TimeSeriesDatum{Labels: labels, Metric: "cpu_temperature", Value: 32.0, Timestamp: 1600000000})
	if err != nil {
		log.Fatalf("could not add: %v", err)
	}

	// Print out the gRPC response.
	log.Printf("Server Response: %s", r.GetMessage())
}
```


### Example 2 - Insert Multiple Rows

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	tspb "github.com/golang/protobuf/ptypes/timestamp"

	pb "github.com/bartmika/tstorage-server/proto"
)

func main() {
  port := 50051

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

	tsd := &pb.TimeSeriesDatum{Labels: labels, Metric: "gpu_temperature, Value: 72, Timestamp: 1600000000}

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
```

### Example 3 - Select

```go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"

	tspb "github.com/golang/protobuf/ptypes/timestamp"

	pb "github.com/bartmika/tstorage-server/proto"
)

func main() {
    port := 50051

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
		Seconds: 1600000000,
		Nanos:   0,
	}
	ets := &tspb.Timestamp{
		Seconds: 1600000009,
		Nanos:   0,
	}

	// Generate our labels.
	labels := []*pb.Label{}
	labels = append(labels, &pb.Label{Name: "Source", Value: "Command"})

	// Perform our gRPC request.
	stream, err := client.Select(ctx, &pb.Filter{Labels: labels, Metric: "battery_percent_left", Start: sts, End: ets})
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
```

## License

This application is licensed under the **BSD 3-Clause License**. See [LICENSE](LICENSE) for more information.

## Acknowledgement

This gRPC server is built overtop [`tstorage`](https://github.com/nakabonne/tstorage) which was architected and written by [Ryo Nakao](https://github.com/nakabonne).
