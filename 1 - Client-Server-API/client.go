package main

import (
	"context"
	"desafios_client_server/utils"
	"errors"
	"fmt"
	"os"
	"time"
)

type cotacao struct {
	Bid string `json:"bid"`
}

const serverURL string = "http://localhost:8080"

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	listenerTimeout(ctx)

	cotacao, err := utils.ExecGetRequestWithContext[cotacao](ctx, serverURL+"/cotacao")
	if err != nil {
		utils.ExitIfError(errors.New(err.Error()), "Erro ao buscar a cotação")
	}

	writeString(cotacao.Bid)

	fmt.Printf("Dólar : %s", cotacao.Bid)
}

func listenerTimeout(ctx context.Context) {

	select {
	case <-ctx.Done():
		utils.ShowError(ctx.Err(), "Tempo de execução da chamda do server insuficiente.")
		return
	default:
		return
	}
}

func writeString(value string) {

	f, err := os.Create("cotacao.txt")

	utils.ExitIfError(err, "Erro ao criar o arquivo de cotações")

	_, err = f.WriteString(fmt.Sprintf("Dólar: %s", value))

	utils.ExitIfError(err, "Erro ao escrever o arquivo de cotações")

	f.Close()
}
