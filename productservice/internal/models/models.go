package models

import (
	"database/sql"
	"github.com/tacheshun/golang-rest-api/productservice/internal/client"
	"time"
)

type Product struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Price    float64   `json:"price"`
	Sale     float64   `json:"total_sales"`
	Created  time.Time `json:"-"`
	Modified time.Time `json:"-"`
}

func (p *Product) GetProduct(db *sql.DB) error {
	return db.QueryRow("SELECT product_id, name, price FROM products WHERE product_id=$1",
		p.ID).Scan(&p.ID, &p.Name, &p.Price)
}

func (p *Product) UpdateProduct(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE products SET name=$1, price=$2 WHERE product_id=$3",
			p.Name, p.Price, p.ID)

	return err
}

func (p *Product) DeleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products WHERE product_id=$1", p.ID)

	return err
}

func (p *Product) CreateProduct(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO products(name, price) VALUES($1, $2) RETURNING product_id",
		p.Name, p.Price).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

func (p *Product) GetProducts(db *sql.DB, start, count int) ([]Product, error) {
	var products []Product
	products = make([]Product, 0)

	rows, err := db.Query(
		"SELECT product_id, name, price FROM products ORDER BY product_id LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		productTotalSales, err := client.GetSalesForProductRPC(uint32(p.ID))
		if err != nil {
			return nil, err
		}
		p.Sale = float64(productTotalSales["quantity"].(uint32))

		products = append(products, p)

	}

	return products, nil
}
