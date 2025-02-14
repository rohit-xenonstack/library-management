package main

import (
	"library-management/backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// dsn := "host=localhost user=postgres password=postgres dbname=library port=5433 sslmode=disable TimeZone=Asia/Kolkata"
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	router := gin.Default()
	router.GET("/ping", handler.Ping)
	router.Run(":8082")
}
