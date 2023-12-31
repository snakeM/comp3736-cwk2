package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/charts/test", generateDataset)
	r.POST("/result/new", handleResultData)
	r.GET("/charts", getChartData)
	r.Run()
}
