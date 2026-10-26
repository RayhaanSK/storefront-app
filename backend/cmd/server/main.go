package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"storefrontapp/domain/storefront/store"
	"storefrontapp/platform/database"
	storefrontservice "storefrontapp/services/storefront"
	storefronthandler "storefrontapp/transport/http/storefront"
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

	repository := store.NewPostgresRepository(db)

	service := storefrontservice.NewService(repository, time.Now)

	handler := storefronthandler.NewHandler(service)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	storefronthandler.RegisterRoutes(router, handler)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
