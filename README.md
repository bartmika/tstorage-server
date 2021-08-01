# tstorage-server
Persistent fast time-series data storage server accessible over gRPC.

`tstorage-server` is lightweight local on-disk storage engine server for time-series data accessible over gRPC. Run this server once and share fast time-series data CRUD operations in either local or remote applications.
The purpose of this server is to allow interprocess communication overtop the [`tstorage`](https://github.com/nakabonne/tstorage) package.

## Installation

Get our latest code.

```bash
go install github.com/bartmika/tstorage-server@latest
```

## Usage

To start the server, run the following command in your **terminal**:

```bash
$GOBIN/tstorage-server serve
```

That's it! If everything works, you should see a message saying `gRPC server is running.`.

## Sub-Commands Reference

### ``serve``

**Details:**

```text
Run the gRPC server to allow other services to access the storage application

Usage:
  tstorage-server serve [flags]

Flags:
  -d, --dataPath string                The location to save the database files to. (default "./tsdb")
  -h, --help                           help for serve
  -b, --partitionDurationInHours int   The timestamp range inside partitions. (default 1)
  -p, --port int                       The port to run this server on (default 50051)
  -t, --timestampPrecision string      The precision of timestamps to be used by all operations. Options:  (default "s")
  -w, --writeTimeoutInSeconds int      The timeout to wait when workers are busy (in seconds). (default 30)
```

**Example:**

```bash
$GOBIN/tstorage-server serve -p=50051 -d="./tsdb" -t="s" -b=1 -w=30
```

### ``insert_row``

**Details:**

```text
Connect to the gRPC server and sends a single time-series datum.

Usage:
  tstorage-server insert_row [flags]

Flags:
  -h, --help            help for insert_row
  -m, --metric string   The metric to attach to the TSD.
  -p, --port int        The port of our server. (default 50051)
  -t, --timestamp int   The timestamp to attach to the TSD.
  -v, --value float     The value to attach to the TSD.
```

**Example:**

```bash
$GOBIN/tstorage-server insert_row -p=50051 -m="solar_biodigester_temperature_in_degrees" -v=50 -t=1600000000
```

Developer Notes:
- There also exists a `insert_rows` subcommand but it works exactly as `insert_row` command with the exception that the internal code is using streaming. This is done so programmers can look at the code and see how to use streaming of time-series data.

### ``select``
**Details:**

```text
Connect to the gRPC server and return list of results based on a selection filter.

Usage:
  tstorage-server select [flags]

Flags:
  -e, --end int         The end timestamp to finish our range
  -h, --help            help for select
  -m, --metric string   The metric to filter by
  -p, --port int        The port of our server. (default 50051)
  -s, --start int       The start timestamp to begin our range
```

**Example:**

```bash
$GOBIN/tstorage-server select --port=50051 --metric="bio_reactor_pressure_in_kpa" --start=1600000000 --end=1725946120
```

## How to Access using gRPC

* Example 1 - Insert a Single Row via [*insert_row.go*](https://github.com/bartmika/tstorage-server/blob/master/cmd/insert_row.go).

* Example 2 - Insert Multiple Rows via [*insert_rows.go*](https://github.com/bartmika/tstorage-server/blob/master/cmd/insert_rows.go).

* Example 3 - Select via [*select.go*](https://github.com/bartmika/tstorage-server/blob/master/cmd/select.go).

* Example 4 - Third Party application via [*poller-server*](https://github.com/bartmika/tpoller-server) code repository.

## What is the gRPC service definition?
Please see the [tstorage.proto](https://github.com/bartmika/tstorage-server/blob/master/proto/tstorage.proto) file for more details. Code snippet from that file:


```protobuf
service TStorage {
    rpc InsertRow (TimeSeriesDatum) returns (google.protobuf.Empty) {}
    rpc InsertRows (stream TimeSeriesDatum) returns (google.protobuf.Empty) {}
    rpc Select (Filter) returns (stream DataPoint) {}
}

message DataPoint {
    double value = 3;
    google.protobuf.Timestamp timestamp = 4;
}

message Label {
    string name = 1;
    string value = 2;
}

message TimeSeriesDatum {
    string metric = 1;
    repeated Label labels = 2;
    double value = 3;
    google.protobuf.Timestamp timestamp = 4;
}

message Filter {
    string metric = 1;
    repeated Label labels = 2;
    google.protobuf.Timestamp start = 3;
    google.protobuf.Timestamp end = 4;
}

message SelectResponse {
    repeated DataPoint points = 1;
}
```

## Contributing
### Development
If you'd like to setup the project for development. Here are the installation steps:

1. Go to your development folder.

    ```bash
    cd ~/go/src/github.com/bartmika
    ```

2. Clone the repository.

    ```bash
    git clone https://github.com/bartmika/tstorage-server.git
    cd tstorage-server
    ```

3. Install the package dependencies

    ```bash
    go mod tidy
    ```

4. In your **terminal**, make sure we export our path (if you haven’t done this before) by writing the following:

    ```bash
    export PATH="$PATH:$(go env GOPATH)/bin"
    ```

5. Run the following to generate our new gRPC interface. Please note in your development, if you make any changes to the gRPC service definition then you'll need to rerun the following:

    ```bash
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/tstorage.proto
    ```

6. You are now ready to start the server and begin contributing!

    ```bash
    go run main.go serve
    ```

### Quality Assurance

Found a bug? Need Help? Please create an [issue](https://github.com/bartmika/tpoller-server/issues).


## License

[**ISC License**](LICENSE) © Bartlomiej Mika

## Acknowledgement

This gRPC server is built overtop [`tstorage`](https://github.com/nakabonne/tstorage) which was architected and written by [Ryo Nakao](https://github.com/nakabonne).
