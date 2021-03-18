package client

import (
	"context"
	"fmt"
	pb "github.com/tacheshun/golang-rest-api/salesservice/proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func Send() (interface{}, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSalesClient(conn)

	// Contact the server and print out its response.
	var productId uint32
	productId = 1
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetProductWithHighestSales(ctx, &pb.ProductIdRequest{ProductId: productId})
	if err != nil {
		log.Fatalf("could not fetch: %v", err)
	}
	fmt.Println(r)
	log.Printf("Product wth the highest sales is: %v, %v", r.Product, r.Sale)

	return r, nil
}
