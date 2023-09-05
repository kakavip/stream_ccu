package main

import (
	"github.com/go-chi/chi/v5"
)

type LimitCCUAPI struct {
}

func (ccuApi *LimitCCUAPI) Routers() *chi.Mux {
	router := chi.NewRouter()

	router.HandleFunc("/", limitCCUHandler)
	return router
}
