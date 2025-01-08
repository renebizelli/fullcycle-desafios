package cep

import (
	"context"
	"fmt"
	"strings"

	"github.com/renebizelli/goexpert/desafios/multithreading/internal/utils"
)

type viaCEPResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Uf          string `json:"uf"`
	Localidade  string `json:"localidade"`
}

type ViaCEP struct {
	Name string `default:"ViaCEP"`
	URL  string
	Ctx  context.Context
}

func NewViaCEPService(url string) *ViaCEP {
	return &ViaCEP{
		URL:  url,
		Name: "ViaCEP",
	}
}

func (v *ViaCEP) Searching(ctx context.Context, searchedCEP string, channel chan<- *Response) {

	url := strings.Replace(v.URL, "?", searchedCEP, 1)

	serviceResponse, error := utils.ExecRequestWithContext[viaCEPResponse](ctx, url)

	if error != nil {
		if error.StatusCode == 499 {
			return
		}
		ErrorMessage(v.Name, searchedCEP, error)
		return
	}

	response := Response{
		Cep:          serviceResponse.Cep,
		State:        serviceResponse.Uf,
		City:         serviceResponse.Localidade,
		Neighborhood: serviceResponse.Bairro,
		Street:       strings.TrimSuffix(fmt.Sprintf("%s %s", serviceResponse.Logradouro, serviceResponse.Complemento), " "),
		Source:       v.Name,
	}

	channel <- &response
}
