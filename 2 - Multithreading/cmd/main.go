package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/renebizelli/goexpert/desafios/multithreading/configs"
	"github.com/renebizelli/goexpert/desafios/multithreading/features/cep"
	"github.com/renebizelli/goexpert/desafios/multithreading/internal/utils"
)

// @title GO Expert - Desafio 2 - Multithreading
// @version 1
// @description Este é o desafio 2 do curso de Go Expert
// @host localhost:8080
// @BasePath /
// @schemes http
// @contact.email rene.bizelli@gmail.com
// @contact.name René Bizelli
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {

	configs := configs.LoadConfig(".")

	fmt.Printf("\n\nAPI is running on port: %s\n\n", utils.YellowText(configs.WebServer.Port))

	mux := chi.NewRouter()

	mux.Use()

	mux.Use(middleware.WithValue("configs", configs))

	cep.AddEndpoint(mux)

	http.ListenAndServe(fmt.Sprintf(":%s", configs.WebServer.Port), mux)

}
