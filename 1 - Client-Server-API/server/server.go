package main

import (
	"context"
	"desafios_client_server/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const serverURL string = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
const pathDB string = "../db/desafio1.db"

type CotacaoServer struct{}

type CotacaoResponse struct {
	USDBRL Cotacao
}

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {

	mux := http.NewServeMux()
	mux.Handle("/cotacao", CotacaoServer{})

	http.ListenAndServe(":8080", mux)

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	// defer cancel()

}

func (s CotacaoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	timeoutObserver(ctx)

	response, err := utils.ExecGetRequestWithContext[CotacaoResponse](ctx, serverURL)

	if err != nil && err.StatusCode != 200 {
		statusCode := err.StatusCode
		if err.StatusCode == 0 {
			statusCode = 408
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(err.Error()))
		return
	}

	db := CotacaoDB{pathDB: pathDB}

	db.Save(response.USDBRL.Bid)
	db.ListStored()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response.USDBRL)
}

func timeoutObserver(ctx context.Context) {

	deadline, _ := ctx.Deadline()

	select {
	case <-ctx.Done():
		utils.ShowError(ctx.Err(), "Chamada da API cancelada")
		return
	case <-time.After(time.Duration(time.Until(deadline))):
		fmt.Println("Tempo de execução da chamada da API insuficiente.")
		return
	default:
	}
}
