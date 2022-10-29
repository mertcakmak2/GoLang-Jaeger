package main

import (
	"context"
	"fmt"
	"go-jaeger/client"
	"go-jaeger/model"
	"log"
	"net/http"
	"time"

	config "go-jaeger/config"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

var tracer opentracing.Tracer

func main() {

	trc, closer, _ := config.InitJaeger()
	tracer = trc
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	router := gin.Default()
	router.GET("/city", getCity)
	router.Run(":8080")
}

func getCity(c *gin.Context) {
	span := config.StartSpanFromRequest(tracer, c.Request, "city-service handle /city")
	defer span.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	if !authenticationControl(ctx) {
		c.String(http.StatusUnauthorized, "unauthorization")
		return
	}

	Country, err := client.GetCountry(ctx, "http://localhost:8081/country")
	if err != nil {
		log.Fatalf("Error occurred: %s", err)
	}
	fmt.Println("Country: ", Country)
	city := model.City{Name: "Trabzon", CountryName: Country.Name}
	c.JSON(http.StatusOK, city)
}

func authenticationControl(ctx context.Context) bool {
	span, _ := opentracing.StartSpanFromContext(ctx, "city-service authenticationControl func")
	defer span.Finish()

	time.Sleep(time.Second)
	return true
}
