package main

import (
	"context"
	"fmt"

	"net"

	"cloud.google.com/go/bigquery"
	pb "github.com/tlgevers/mapping-bg/proto"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
)

const (
	projectID = "tgevers-apps"
	port      = ":50051"
)

type server struct {
	pb.UnimplementedFAAAirportDataServer
}

func (s *server) GetAirportData(ctx context.Context, in *pb.RequestCode) (*pb.Airports, error) {
	c := context.Background()
	client, err := bigquery.NewClient(c, projectID)
	if err != nil {
		fmt.Println(err)
	}
	q := client.Query(`
	    SELECT faa_identifier,name,latitude,longitude
		FROM ` + "`bigquery-public-data.faa.us_airports`" + `
		LIMIT 10
	`)
	it, err := q.Read(ctx)
	if err != nil {
		fmt.Println(err) // TODO: Handle error.
	}
	var ap pb.Airport
	var aps []*pb.Airport
	for {
		err := it.Next(&ap)
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		aps = append(aps, &ap)
	}
	fmt.Println(aps)
	fmt.Printf("Received: %v", in.GetCode())
	// var a pb.Airport
	// for _, ap := range(values) {
	// 	err = json.Unmarshal(ap, &a)
	// 	aps = append(aps, )
	// }
	return &pb.Airports{
		Airports: aps,
	}, nil
}

func testGetData() {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		fmt.Println(err)
	}
	q := client.Query(`
	    SELECT faa_identifier,name,latitude,longitude
		FROM ` + "`bigquery-public-data.faa.us_airports`" + `
		LIMIT 10
	`)
	it, err := q.Read(ctx)
	if err != nil {
		fmt.Println(err) // TODO: Handle error.
	}
	var values []pb.Airport
	var ap pb.Airport
	for {
		err := it.Next(&ap)
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(values)
		values = append(values, ap)

	}
	fmt.Println(values)
}

func main() {
	fmt.Println("vim-go")
	// testGetData()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFAAAirportDataServer(s, &server{})
	fmt.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
	}
}
