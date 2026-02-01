package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/willianbsanches13/pack-calculator/internal/handler"
	"github.com/willianbsanches13/pack-calculator/internal/storage"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("GIN_MODE")
	if env == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	store := storage.NewMemoryStorage()
	h := handler.New(store)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(corsMiddleware())

	h.RegisterRoutes(r)

	log.Printf("Pack Calculator API running on http://localhost:%s", port)
	log.Printf("Default pack sizes: %v", store.GetPackSizes())

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server failed:", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
