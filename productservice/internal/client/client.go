package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	pb "github.com/tacheshun/golang-rest-api/salesservice/proto"
	"github.com/uber/jaeger-client-go/config"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func InitTracer(serviceName, host string) (opentracing.Tracer, error) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  host,
		},
	}

	tracer, _, err := cfg.New(serviceName)
	if err != nil {
		return nil, fmt.Errorf("new tracer error: %v", err)
	}
	return tracer, nil
}

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
	log.Printf("Product wth the highest sales is: %v, %v", r.Product, r.TotalSales)

	return r, nil
}
