package main

import (
	"net/http"
	"time"

	config "country-service/config"
	model "country-service/model"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func main() {
	_, closer, _ := config.InitJaeger()
	defer closer.Close()
	opentracing.SetGlobalTracer(config.Tracer)

	router := gin.Default()
	router.GET("/country", getCountry)
	router.Run(":8081")
}

func getCountry(c *gin.Context) {
	span := config.StartSpanFromRequest(config.Tracer, c.Request.Header, "country-service handle /country")
	defer span.Finish()
	time.Sleep(500 * time.Millisecond)
	isThereCountry(c.Request.Header)
	country := model.Country{Name: "Turkey"}
	c.JSON(http.StatusOK, country)
}

func isThereCountry(h http.Header) bool {
	span := config.StartSpanFromRequest(config.Tracer, h, "country-service isThereCountry func")
	defer span.Finish()
	time.Sleep(200 * time.Millisecond)
	return true
}
