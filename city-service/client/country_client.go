package client

import (
	"context"
	"encoding/json"
	"go-jaeger/config"
	"go-jaeger/model"
	"go-jaeger/request"
	"net/http"

	"github.com/opentracing/opentracing-go"
)

func GetCountry(ctx context.Context, hostPort string) (model.Country, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "city-service client/GetCountry func")
	defer span.Finish()

	req, err := http.NewRequest("GET", hostPort, nil)
	if err != nil {
		return model.Country{}, err
	}

	if err := config.Inject(span, req); err != nil {
		return model.Country{}, err
	}

	country := model.Country{}
	data, _ := request.Do(req)
	json.Unmarshal(data, &country)
	return country, nil
}
