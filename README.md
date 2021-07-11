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

## Sub-Commands

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

## License

This application is licensed under the **BSD 3-Clause License**. See [LICENSE](LICENSE) for more information.
