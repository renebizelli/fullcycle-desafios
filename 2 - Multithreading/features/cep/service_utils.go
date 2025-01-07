package cep

import (
	"fmt"

	"github.com/renebizelli/goexpert/desafios/multithreading/internal/utils"
)

func ErrorMessage(serviceName string, cep string, err error) {
	fmt.Printf("\n%s \ncep: %s.\n %s", utils.RedText("error "+serviceName), cep, err.Error())
}
