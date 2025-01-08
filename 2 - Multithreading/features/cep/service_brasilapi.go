package cep

import (
	"context"
	"fmt"

	"github.com/renebizelli/goexpert/desafios/multithreading/internal/utils"
)

type brasilAPIResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

type BrasilApi struct {
	Name string
	URL  string
	Ctx  context.Context
}

func NewBrasilAPIService(url string) *BrasilApi {
	return &BrasilApi{
		URL:  url,
		Name: "BrasilAPI",
	}
}

func (b *BrasilApi) Searching(ctx context.Context, searchedCEP string, channel chan<- *Response) {

	url := fmt.Sprintf("%s%s", b.URL, searchedCEP)

	brasilAPIResponse, error := utils.ExecRequestWithContext[brasilAPIResponse](ctx, url)

	if error != nil {
		if error.StatusCode == 499 {
			return
		}
		ErrorMessage(b.Name, searchedCEP, error)
		return
	}

	response := Response{
		Cep:          brasilAPIResponse.Cep,
		State:        brasilAPIResponse.State,
		City:         brasilAPIResponse.City,
		Neighborhood: brasilAPIResponse.Neighborhood,
		Street:       brasilAPIResponse.Street,
		Source:       b.Name,
	}

	channel <- &response
}
