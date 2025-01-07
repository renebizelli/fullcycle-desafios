package cep

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/renebizelli/goexpert/desafios/multithreading/internal/utils"
)

type serviceResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

type BrasilApi struct {
	Name    string
	URL     string
	Timeout int
}

func NewBrasilAPIService(url string, timeout int) *BrasilApi {
	return &BrasilApi{
		URL:     url,
		Timeout: timeout,
		Name:    "BrasilAPI",
	}
}

func (b *BrasilApi) Searching(searchedCEP string) (*Response, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(b.Timeout))
	defer cancel()

	url := fmt.Sprintf("%s%s", b.URL, searchedCEP)

	serviceResponse, error := utils.ExecRequestWithContext[serviceResponse](ctx, url)

	if error != nil {
		return nil, errors.New(ErrorMessage(b.Name, searchedCEP, error))
	}

	response := Response{
		Cep:          serviceResponse.Cep,
		State:        serviceResponse.State,
		City:         serviceResponse.City,
		Neighborhood: serviceResponse.Neighborhood,
		Street:       serviceResponse.Street,
		Source:       b.Name,
	}

	return &response, nil

}
