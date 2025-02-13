package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Defining an interface for the custom Env type that wraps the DB connection handle
// Implemented using the idea from Alex Edwards blog on "Organising Database Access in Go"
type Env struct {
	library1 interface{}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "service up",
		"timestamp": time.Now(),
	})
}

func main() {
	// dsn := "host=localhost user=postgres password=postgres dbname=library port=5433 sslmode=disable TimeZone=Asia/Kolkata"
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	router := gin.Default()
	router.GET("/check", healthCheck)
	router.Run(":8082")
}
