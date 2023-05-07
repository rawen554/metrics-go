package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	storage map[string]float64
}

const (
	gauge   = "gauge"
	counter = "counter"
)

func updateMetric(memStorage *MemStorage) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		reqURI := strings.Split(req.URL.Path, "/")
		if len(reqURI) != 5 {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		metricType := reqURI[2]
		metricName := reqURI[3]
		if metricName == "" {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		memStorage := *memStorage
		storage := memStorage.storage

		switch metricType {
		case gauge:
			metricValue, err := strconv.ParseFloat(reqURI[4], 64)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			storage[metricName] = metricValue

		case counter:
			metricValue, err := strconv.ParseInt(reqURI[4], 10, 64)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			currentValue := int64(storage[metricName])
			storage[metricName] = float64(currentValue + metricValue)
		default:
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
	}
}

func main() {
	memStorage := MemStorage{storage: map[string]float64{}}
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", updateMetric(&memStorage))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
