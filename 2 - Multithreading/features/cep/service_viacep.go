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
	Name    string `default:"ViaCEP"`
	URL     string
	Timeout int
	Ctx     context.Context
}

func NewViaCEPService(url string, ctx context.Context) *ViaCEP {
	return &ViaCEP{
		URL:  url,
		Ctx:  ctx,
		Name: "ViaCEP",
	}
}

func (v *ViaCEP) Searching(searchedCEP string, channel chan<- *Response) {

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(v.Timeout))
	// defer cancel()

	url := strings.Replace(v.URL, "?", searchedCEP, 1)

	serviceResponse, error := utils.ExecRequestWithContext[viaCEPResponse](v.Ctx, url)

	if error != nil {
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
