package cep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/renebizelli/goexpert/desafios/multithreading/configs"
	"github.com/renebizelli/goexpert/desafios/multithreading/internal/dtos"
	"github.com/renebizelli/goexpert/desafios/multithreading/internal/utils"
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

	e := cepValidate(searchedCEP)

	if e != nil {
		fmt.Printf("\n Invalid cep %s\n\n", utils.RedText(searchedCEP))
		w.WriteHeader(400)
		w.Write([]byte(e.Error()))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(config.Services.Timeout))
	defer cancel()

	brasilApi := NewBrasilAPIService(config.Services.BrasilApiUrl)
	viaCEP := NewViaCEPService(config.Services.ViacepUrl)

	services := []SearchingInfo{brasilApi, viaCEP}

	ch := make(chan *Response)

	for _, service := range services {
		go service.Searching(ctx, searchedCEP, ch)
	}

	select {
	case response := <-ch:
		ctx.Done()

		printResponse(response)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response)

	case <-time.After(3 * time.Second):
		fmt.Printf("\nRequest %s for CEP %s", utils.RedText("timeout"), utils.CyanText(searchedCEP))
		w.WriteHeader(408)
	}

}

func cepValidate(searchedCEP string) error {

	if len(searchedCEP) != 8 {
		return errors.New("invalid cep")
	}

	re := regexp.MustCompile("[0-9]+")
	searchedCEP = strings.Join(re.FindAllString(searchedCEP, -1)[:], "")

	if len(searchedCEP) != 8 {
		return errors.New("invalid cep")
	}

	return nil

}

func printResponse(response *Response) {
	fmt.Println("\n--------------------------------------------------")
	fmt.Printf("Response from %s", utils.YellowText(response.Source))
	fmt.Println("\n--------------------------------------------------")
	fmt.Printf("CEP..........: %s\n", utils.CyanText(response.Cep))
	fmt.Printf("State........: %s\n", utils.CyanText(response.State))
	fmt.Printf("City.........: %s\n", utils.CyanText(response.City))
	fmt.Printf("Neighborhood.: %s\n", utils.CyanText(response.Neighborhood))
	fmt.Printf("Street.......: %s\n", utils.CyanText(response.Street))
	fmt.Println("--------------------------------------------------")
}
