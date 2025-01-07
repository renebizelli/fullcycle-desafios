package cep

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/renebizelli/goexpert/desafios/multithreading/configs"
	"github.com/renebizelli/goexpert/desafios/multithreading/internal/dtos"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	config := r.Context().Value("configs").(configs.Config)

	if (config == configs.Config{}) {
		w.WriteHeader(500)
		e := dtos.RequestError{Message: "Config is empty", StatusCode: 500}
		json.NewEncoder(w).Encode(e)
		return
	}

	searchedCEP := chi.URLParam(r, "cep")

	viaCEP := NewViaCEPService(config.Services.ViacepUrl, config.Services.Timeout)
	brasilApi := NewBrasilAPIService(config.Services.BrasilApiUrl, config.Services.Timeout)

	services := []SearchingInfo{viaCEP, brasilApi}

	w.Header().Set("Content-Type", "application/json")

	//wg := sync.WaitGroup{}

	for _, service := range services {

		response, error := service.Searching(searchedCEP)

		if error != nil {
			fmt.Println(error.Error())
		} else {
			json.NewEncoder(w).Encode(response)
		}
	}

	w.WriteHeader(200)
}
