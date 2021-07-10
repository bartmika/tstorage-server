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


## License

This application is licensed under the **BSD 3-Clause License**. See [LICENSE](LICENSE) for more information.
