package main

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	pb "github.com/tacheshun/golang-rest-api/salesservice/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
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
	grpcServer := grpc.NewServer()
	pb.RegisterSalesServer(grpcServer, &server{
		repo: &app,
	})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type SalesApp struct {
	DB *sql.DB
}

type Sale struct {
	SalesID   int `json:"sale_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
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
	repo *SalesApp
}

func (srv *server) GetProductWithHighestSales(ctx context.Context, in *pb.ProductIdRequest) (*pb.ProductWithSales, error) {
	result, err := srv.repo.GetProductPlusHighestSales(in.GetProductId())
	if err != nil {
		return nil, err
	}
	out := &pb.ProductWithSales{}

	out.Product = uint32(result.ProductID)
	out.TotalSales = uint32(result.Quantity)

	return out, nil
}

func (s *SalesApp) GetProductPlusHighestSales(_ uint32) (*ProductSales, error) {
	var productSales ProductSales
	_ = s.DB.QueryRow(
		"SELECT product_id, SUM(quantity) as quantity from sales GROUP BY product_id ORDER BY quantity DESC LIMIT 1").Scan(&productSales.ProductID, &productSales.Quantity)

	return &productSales, nil
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
