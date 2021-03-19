package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	opentracing "github.com/opentracing/opentracing-go"
	pb "github.com/tacheshun/golang-rest-api/salesservice/proto"
	"github.com/uber/jaeger-client-go/config"
)

const (
	address     = ":50051"
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
	tracer, err := InitTracer("product", "127.0.0.1:16686")
	if err != nil {
		panic("cannot start tracer")
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(
		otgrpc.OpenTracingClientInterceptor(tracer)), grpc.WithStreamInterceptor(
		otgrpc.OpenTracingStreamClientInterceptor(tracer)), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSalesClient(conn)

	// Contact the server and print out its response.
	var productId uint32
	productId = 2
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetProductWithHighestSales(ctx, &pb.ProductIdRequest{ProductId: productId})
	if err != nil {
		log.Fatalf("could not fetch: %v", err)
	}

	return r, nil
}

func GetSalesForProductRPC(productId uint32) (map[string]interface{}, error) {
	tracer, err := InitTracer("product", "127.0.0.1:16686")
	if err != nil {
		panic("cannot start tracer")
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(
		otgrpc.OpenTracingClientInterceptor(tracer)), grpc.WithStreamInterceptor(
		otgrpc.OpenTracingStreamClientInterceptor(tracer)), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSalesClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.GetSalesForProduct(ctx, &pb.ProductIdRequest{ProductId: productId})
	if err != nil {
		log.Fatalf("could not fetch: %v", err)
	}

	return map[string]interface{}{"productId":res.ProductId, "quantity":res.Quantity}, nil
}
