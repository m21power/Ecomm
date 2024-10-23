package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/m21power/ecomm/ecomm-api/server"
	"github.com/m21power/ecomm/ecomm-api/storer"
)

type handler struct {
	ctx    context.Context
	server *server.Server
}

func NewHandler(server *server.Server) *handler {
	return &handler{ctx: context.Background(), server: server}
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p ProductReq
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}
	createdProduct, err := h.server.CreateProduct(h.ctx, toStorerProduct(p))
	if err != nil {
		http.Error(w, "error creating product", http.StatusInternalServerError)
		return
	}
	res := toProductRes(createdProduct)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)

}

func (h *handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing id", http.StatusBadRequest)
		return
	}
	product, err := h.server.GetProduct(h.ctx, i)
	if err != nil {
		http.Error(w, "error getting product", http.StatusInternalServerError)
		return
	}
	res := toProductRes(product)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.server.ListProducts(h.ctx)
	if err != nil {
		http.Error(w, "error listing products", http.StatusInternalServerError)
		return
	}
	var res []*ProductRes
	for _, p := range products {
		res = append(res, toProductRes(&p))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing id", http.StatusBadRequest)
		return
	}
	var p ProductReq
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}
	product, err := h.server.GetProduct(h.ctx, i)
	if err != nil {
		http.Error(w, "error getting product", http.StatusInternalServerError)
		return
	}
	// now it is a time to update the product
	toPatchProduct(product, p)
	updatedProduct, err := h.server.UpdateProduct(h.ctx, product)
	if err != nil {
		http.Error(w, "error updating product", http.StatusInternalServerError)
		return
	}
	res := toProductRes(updatedProduct)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

}

func (h *handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing id", http.StatusBadRequest)
		return
	}
	err = h.server.DeleteProduct(h.ctx, i)
	if err != nil {
		http.Error(w, "error deleting product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func toStorerProduct(p ProductReq) *storer.Product {
	return &storer.Product{
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Rating:       p.Rating,
		NumReviews:   p.NumReviews,
		Price:        p.Price,
		CountInStock: p.CountInStock,
		CreatedAt:    time.Now(),
	}
}

func toProductRes(p *storer.Product) *ProductRes {
	return &ProductRes{
		ID:           p.ID,
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Rating:       p.Rating,
		NumReviews:   p.NumReviews,
		Price:        p.Price,
		CountInStock: p.CountInStock,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

func toPatchProduct(product *storer.Product, p ProductReq) {
	if p.Name != "" {
		product.Name = p.Name
	}
	if p.Image != "" {
		product.Image = p.Image
	}
	if p.Category != "" {
		product.Category = p.Category
	}
	if p.Description != "" {
		product.Description = p.Description
	}
	if p.Rating != 0 {
		product.Rating = p.Rating
	}
	if p.NumReviews != 0 {
		product.NumReviews = p.NumReviews
	}
	if p.Price != 0 {
		product.Price = p.Price
	}
	if p.CountInStock != 0 {
		product.CountInStock = p.CountInStock
	}
	product.UpdatedAt = toTimePtr(time.Now())
}

func toTimePtr(t time.Time) *time.Time {
	return &t
}
