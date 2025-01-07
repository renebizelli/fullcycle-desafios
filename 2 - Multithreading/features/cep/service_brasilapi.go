package cep

import (
	"context"
	"fmt"

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
	Ctx     context.Context
}

func NewBrasilAPIService(url string, ctx context.Context) *BrasilApi {
	return &BrasilApi{
		URL:  url,
		Ctx:  ctx,
		Name: "BrasilAPI",
	}
}

func (b *BrasilApi) Searching(searchedCEP string, channel chan<- *Response) {

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(b.Timeout))
	// defer cancel()

	url := fmt.Sprintf("%s%s", b.URL, searchedCEP)

	serviceResponse, error := utils.ExecRequestWithContext[serviceResponse](b.Ctx, url)

	if error != nil {
		ErrorMessage(b.Name, searchedCEP, error)
		return
	}

	response := Response{
		Cep:          serviceResponse.Cep,
		State:        serviceResponse.State,
		City:         serviceResponse.City,
		Neighborhood: serviceResponse.Neighborhood,
		Street:       serviceResponse.Street,
		Source:       b.Name,
	}

	channel <- &response
}
