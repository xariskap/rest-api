package main

import (
	"context"
	"fmt"
	"net/http"

	"rest/db"

	"github.com/gin-gonic/gin"
)

var dbConn = db.GetConnection()
var dbName = "Simpler"

type Product struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

// Get all products
func getAllProducts(db *db.Database) gin.HandlerFunc {

	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, price, quantity FROM products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var products []Product
		for rows.Next() {
			var p Product
			err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			products = append(products, p)
		}

		c.JSON(http.StatusOK, products)
	}
}

// Add a product
func addProduct(db *db.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := db.ExecQuery("INSERT INTO products (name, price, quantity) VALUES ($1, $2, $3)", product.Name, product.Price, product.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, product)
	}
}

func main() {
	database := db.USE(dbName, dbConn)
	defer dbConn.Close(context.Background())

	r := gin.Default()

	r.GET("/products", getAllProducts(database))
	r.POST("/products", addProduct(database))

	fmt.Println("Server is running on http://localhost:8888")
	if err := r.Run(":8888"); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
