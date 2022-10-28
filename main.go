package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func main() {

	tracer, closer, err := InitJaeger()
	if err != nil {
		fmt.Printf("error init jaeger %v", err)
	} else {
		opentracing.SetGlobalTracer(tracer)
		defer closer.Close()
	}

	router := gin.Default()
	router.GET("/cities", getCities)
	router.Run(":8080")
}

func InitJaeger() (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeRateLimiting,
			Param: 100, // 100 traces per second
		},
	}

	tracer, closer, err := cfg.New("go-jaeger-gin-service")
	return tracer, closer, err
}

func getCities(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "Handle /cities")
	defer span.Finish()

	if !authenticationControl(ctx) {
		c.String(http.StatusUnauthorized, "unauthorization")
		return
	}

	cities, err := getCityWithCountryName(ctx)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, cities)
}

func getCityWithCountryName(ctx context.Context) ([]City, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "getCityWithCountryName")
	defer span.Finish()

	cities := []City{
		{Name: "Eskişehir", CountryName: "TURKEY"},
		{Name: "Trabzon", CountryName: "TURKEY"},
		{Name: "Berlin", CountryName: "GERMANY"},
	}

	time.Sleep(500 * time.Millisecond)
	return cities, nil
}

func authenticationControl(ctx context.Context) bool {
	span, _ := opentracing.StartSpanFromContext(ctx, "authenticationControl")
	defer span.Finish()

	time.Sleep(time.Second)
	return true
}

type City struct {
	Name        string `json:"name"`
	CountryName string `json:"country"`
}
