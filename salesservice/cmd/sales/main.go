package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	pb "github.com/tacheshun/golang-rest-api/salesservice/proto"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	port = ":50051"
)

func main() {
	app := SalesApp{}
	app.Initialize()
	log.Println("Starting server Sales localhost on port 50051...")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	tracer, err := InitTracer("Sales", "127.0.0.1:16686")
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		otgrpc.OpenTracingServerInterceptor(tracer)),
		grpc.StreamInterceptor(
			otgrpc.OpenTracingStreamServerInterceptor(tracer)))
	pb.RegisterSalesServer(grpcServer, &server{
		repo:   &app,
		tracer: tracer,
	})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type SalesApp struct {
	DB *sql.DB
}

type Sale struct {
	SalesID   int       `json:"sale_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Created   time.Time `json:"created"`
}

type ProductSales struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

func (s *SalesApp) Initialize() {
	var err error
	s.DB, err = sql.Open("postgres", "user=marius password=magic dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

type server struct {
	pb.UnimplementedSalesServer
	repo   *SalesApp
	tracer opentracing.Tracer
}

func (srv *server) GetProductWithHighestSales(ctx context.Context, in *pb.ProductIdRequest) (*pb.ProductWithSales, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := srv.tracer.StartSpan("GetProductWithHighestSales", opentracing.ChildOf(span.Context()))
		span.SetTag("param.location", in.ProductId)
		ext.SpanKindRPCClient.Set(span)
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	result, err := srv.repo.GetProductPlusHighestSales(in.GetProductId())
	if err != nil {
		return nil, err
	}
	out := &pb.ProductWithSales{}

	out.Product = uint32(result.ProductID)
	out.TotalSales = uint32(result.Quantity)

	return out, nil
}

func (srv *server) GetSalesForProduct(ctx context.Context, in *pb.ProductIdRequest) (*pb.Sale, error) {
	result, err := srv.repo.GetSalesForProductID(in.GetProductId())
	if err != nil {
		return nil, err
	}

	out := &pb.Sale{}
	//
	out.ProductId = uint32(result.ProductID)
	out.Quantity = uint32(result.Quantity)
	return out, nil
}

func (s *SalesApp) GetProductPlusHighestSales(_ uint32) (*ProductSales, error) {
	var productSales ProductSales
	_ = s.DB.QueryRow(
		"SELECT product_id, SUM(quantity) as quantity from sales GROUP BY product_id ORDER BY quantity DESC LIMIT 1").Scan(&productSales.ProductID, &productSales.Quantity)

	return &productSales, nil
}

func (s *SalesApp) GetSalesForProductID(productId uint32) (*Sale, error) {
	var sale Sale
	_ = s.DB.QueryRow(
		"SELECT product_id, sum(quantity) as total_sales from sales WHERE product_id=$1 group by product_id", productId).Scan(&sale.ProductID, &sale.Quantity)

	return &sale, nil
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

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
