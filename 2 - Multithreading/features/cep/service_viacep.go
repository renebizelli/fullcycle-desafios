package cep

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

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
}

func NewViaCEPService(url string, timeout int) *ViaCEP {
	return &ViaCEP{
		URL:     url,
		Timeout: timeout,
		Name:    "ViaCEP",
	}
}

func (v *ViaCEP) Searching(searchedCEP string) (*Response, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(v.Timeout))
	defer cancel()

	url := strings.Replace(v.URL, "?", searchedCEP, 1)

	serviceResponse, error := utils.ExecRequestWithContext[viaCEPResponse](ctx, url)

	if error != nil {
		return nil, errors.New(ErrorMessage(v.Name, searchedCEP, error))
	}

	response := Response{
		Cep:          serviceResponse.Cep,
		State:        serviceResponse.Uf,
		City:         serviceResponse.Localidade,
		Neighborhood: serviceResponse.Bairro,
		Street:       strings.TrimSuffix(fmt.Sprintf("%s %s", serviceResponse.Logradouro, serviceResponse.Complemento), " "),
		Source:       v.Name,
	}

	return &response, nil
}
