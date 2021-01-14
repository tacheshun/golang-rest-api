package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
)

func main() {
	s := SalesApp{}
	s.Initialize()
	log.Println("Starting server Sales localhost on port 8001...")
	s.Run("localhost:8001")
}

type SalesApp struct {
	Router *mux.Router
	DB     *sql.DB
}

func (s *SalesApp) initializeRoutes() {
	s.Router.HandleFunc("/sales/{saleId}", s.GetSale).Methods("GET")
	s.Router.HandleFunc("/sales/product/{productId}", s.GetSaleForProduct).Methods("GET")
}

func (s *SalesApp) Initialize() {
	var err error
	s.DB, err = sql.Open("postgres", "user=marius password=magic dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	s.Router = mux.NewRouter()
	s.initializeRoutes()
}

func (s *SalesApp) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, s.Router))
}

func (s *SalesApp) GetSale(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	saleId, err := strconv.Atoi(vars["saleId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid saleId ID")
		return
	}

	sale := Sale{SalesID: saleId}
	if err := sale.getSaleByID(s.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Sale not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, sale)
}
func (s *SalesApp) GetSaleForProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["productId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid productId ID")
		return
	}

	totalSales := TotalSales{ProductID: productId}
	if err := totalSales.getTotalSalesForProduct(s.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, totalSales)
}

type Sale struct {
	SalesID	   int  `json:"sale_id"`
	ProductID  int  `json:"product_id"`
	Quantity int    `json:"quantity"`
}

type TotalSales struct {
	ProductID int `json:"product_id" db:"product_id"`
	SalesID   int `json:"total_sales,omitempty" db:"sales_id"`
	Total  int `json:"total_quantity_sold,omitempty" db:"quantity"`
}

func (sale *Sale) getSaleByID(db *sql.DB) error {
	return db.QueryRow("SELECT sale_id, product_id, quantity FROM sales WHERE sale_id=$1",
		sale.SalesID).Scan(&sale.SalesID, &sale.ProductID, &sale.Quantity)
}

func (ts *TotalSales) getTotalSalesForProduct(db *sql.DB) error {
	return db.QueryRow("SELECT product_id, count(sale_id), SUM(quantity) as TotalQty from sales where product_id=$1 GROUP BY product_id",
		ts.ProductID).Scan(&ts.ProductID, &ts.SalesID, &ts.Total)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
