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
