package main

import (
	"context"
	"database/sql"
	"desafios_client_server/utils"
	"fmt"
	"strconv"
	"time"
)

type CotacaoDB struct {
	pathDB string
}

func (c CotacaoDB) InitDB() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", c.pathDB)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Cotacao (Value decimal(3, 4))")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (s CotacaoDB) Save(bid string) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	timeoutObserverDB(ctx)

	db, err := s.InitDB()
	utils.ShowError(err, "Erro ao abrir o db")
	defer db.Close()

	value, e := strconv.ParseFloat(bid, 32)
	utils.ShowError(e, "Erro ao converter o Bid para float")

	_, err = db.ExecContext(ctx, "INSERT INTO Cotacao (value) VALUES (?)", value)
	utils.ShowError(err, "Erro ao executer o commnad")

}

func (s CotacaoDB) ListStored() {

	db, err := s.InitDB()
	utils.ShowError(err, "Erro ao abrir o db")
	defer db.Close()

	query, err := db.Query("SELECT * FROM Cotacao")
	utils.ShowError(err, "Erro ao criar query")

	for query.Next() {
		var value float32
		query.Scan(&value)
		fmt.Printf("\nValue: %v", value)
	}

	e := query.Err()
	utils.ShowError(e, "Erro ao listar as cotações")

	defer query.Close()
}

func timeoutObserverDB(ctx context.Context) {

	deadline, _ := ctx.Deadline()

	select {
	case <-ctx.Done():
		utils.ShowError(ctx.Err(), "Operação no DB cancelada")
		return
	case <-time.After(time.Duration(time.Until(deadline))):
		fmt.Println("Tempo de execução da operação no BD insuficiente.")
		return
	default:
	}
}
