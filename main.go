package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	pb "github.com/tlgevers/mapping-bg/proto/bq/bq"
	"google.golang.org/api/iterator"
)

const (
	projectID = "tgevers-apps"
	port      = ":50051"
)

type server struct {
	pb.UnimplementedFAAAirportDataServer
}

func (s *server) GetAirportData(ctx context.Context, in *pb.RequestCode) *pb.Airports {
	fmt.Printf("Received: %v", in.GetName())
	return &pb.Airports{}
}

func main() {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		fmt.Println(err)
	}
	q := client.Query(`
		SELECT year, SUM(number) as num
		FROM ` + "`bigquery-public-data.usa_names.usa_1910_2013`" + `
		WHERE name = "William"
		GROUP BY year
		ORDER BY year
	`)
	it, err := q.Read(ctx)
	if err != nil {
		fmt.Println(err) // TODO: Handle error.
	}
	for {
		var values []bigquery.Value
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(values)
	}

	fmt.Println("vim-go")
}
