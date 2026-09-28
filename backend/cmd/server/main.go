package main

import (
	"context"
	"log"
	"net/http"
	"storefrontapp/platform/database"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.OpenPostgres(ctx, database.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "storefront",
		Password: "storefront",
		Name:     "storefrontapp",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
