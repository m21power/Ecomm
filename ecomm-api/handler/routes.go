package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

var r *chi.Mux

func RegisterRoutes(handler *handler) *chi.Mux {
	r = chi.NewRouter()
	r.Route("/products", func(r chi.Router) {
		r.Post("/", handler.CreateProduct)
		r.Get("/", handler.ListProducts)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetProduct)
			r.Patch("/", handler.UpdateProduct)
			r.Delete("/", handler.DeleteProduct)
		})
	})
	return r
}

func Start(listenAddr string) error {
	return http.ListenAndServe(listenAddr, r)
}
