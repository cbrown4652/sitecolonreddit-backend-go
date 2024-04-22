package main

import (
	"net/http"
	"os"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	router.PUT("/search", searchPUT)

	port := os.Getenv("PORT")
	router.Run(port)
}
