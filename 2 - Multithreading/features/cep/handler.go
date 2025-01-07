package cep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(config.Services.Timeout))
	defer cancel()

	brasilApi := NewBrasilAPIService(config.Services.BrasilApiUrl, ctx)
	viaCEP := NewViaCEPService(config.Services.ViacepUrl, ctx)

	services := []SearchingInfo{brasilApi, viaCEP}

	ch := make(chan *Response)

	for _, service := range services {
		go service.Searching(searchedCEP, ch)
	}

	select {
	case response := <-ch:
		ctx.Done()
		fmt.Printf("Received from %s: %v\n", response.Cep, response.Source)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(200)

	case <-time.After(3 * time.Second):
		fmt.Println("timeout")
		w.WriteHeader(400)

	}

}
