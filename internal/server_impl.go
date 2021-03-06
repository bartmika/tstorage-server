package internal

import (
	"context"
	"io"

	"github.com/golang/protobuf/ptypes/empty"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/nakabonne/tstorage"

	pb "github.com/bartmika/tstorage-server/proto"
)

type TStorageServerImpl struct {
	storage tstorage.Storage
	pb.TStorageServer
}

func (s *TStorageServerImpl) InsertRow(ctx context.Context, in *pb.TimeSeriesDatum) (*empty.Empty, error) {
	// Generate our labels, if there are any.
	labels := []tstorage.Label{}
	for _, label := range in.Labels {
		labels = append(labels, tstorage.Label{Name: label.Name, Value: label.Value})
	}

	// Generate our datapoint.
	dataPoint := tstorage.DataPoint{Timestamp: in.Timestamp.Seconds, Value: in.Value}

	err := s.storage.InsertRows([]tstorage.Row{
		{
			Metric:    in.Metric,
			Labels:    labels,
			DataPoint: dataPoint,
		},
	})
	return &empty.Empty{}, err
}

func (s *TStorageServerImpl) InsertRows(stream pb.TStorage_InsertRowsServer) error {
	// DEVELOPERS NOTE:
	// If you don't understand how server side streaming works using gRPC then
	// please visit the documentation to get an understanding:
	// https://grpc.io/docs/languages/go/basics/#server-side-streaming-rpc-1

	// Wait and receieve the stream from the client.
	for {
		datum, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&empty.Empty{})
		}
		if err != nil {
			return err
		}

		// Generate our labels, if there are any.
		labels := []tstorage.Label{}
		for _, label := range datum.Labels {
			labels = append(labels, tstorage.Label{Name: label.Name, Value: label.Value})
		}

		// Generate our datapoint.
		dataPoint := tstorage.DataPoint{Timestamp: datum.Timestamp.Seconds, Value: datum.Value}

		err = s.storage.InsertRows([]tstorage.Row{
			{
				Metric:    datum.Metric,
				Labels:    labels,
				DataPoint: dataPoint,
			},
		})
	}

	return nil
}

func (s *TStorageServerImpl) Select(in *pb.Filter, stream pb.TStorage_SelectServer) error {
	// Generate our labels, if there are any.
	labels := []tstorage.Label{}
	for _, label := range in.Labels {
		labels = append(labels, tstorage.Label{Name: label.Name, Value: label.Value})
	}

	points, err := s.storage.Select(in.Metric, labels, in.Start.Seconds, in.End.Seconds)
	if err != nil {
		return err
	}

	for _, point := range points {
		ts := &tspb.Timestamp{
			Seconds: point.Timestamp,
			Nanos:   0,
		}
		dataPoint := &pb.DataPoint{Value: point.Value, Timestamp: ts}
		if err := stream.Send(dataPoint); err != nil {
			return err
		}
	}

	return nil
}
