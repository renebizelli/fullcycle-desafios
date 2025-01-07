package cep

import (
	"github.com/go-chi/chi/v5"
)

func AddEndpoint(mux *chi.Mux) {

	mux.Route("/cep", func(c chi.Router) {
		c.Get("/{cep}", Handler)
	})

}
