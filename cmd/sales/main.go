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
	s.Router.HandleFunc("/sales/{product:[0-9]+}", s.GetSaleForProduct).Methods("GET")
}

func (s *SalesApp) Initialize() {
	var err error
	s.DB, err = sql.Open("postgres", "user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	s.Router = mux.NewRouter()
	s.initializeRoutes()
}

func (s *SalesApp) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, s.Router))
}

func (s *SalesApp) GetSaleForProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["product"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	sale := Sale{ProductID: productID}
	if err := sale.getSaleForProduct(s.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product or Sale not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, sale)
}

type Sale struct {
	ProductID  int  `json:"product_id"`
	TotalSales int  `json:"total_sales"`
}

func (sale *Sale) getSaleForProduct(db *sql.DB) error {
	return db.QueryRow("SELECT product_id, total_sales FROM sales WHERE product_id=$1",
		sale.ProductID).Scan(&sale.ProductID, &sale.TotalSales)
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
