package server

import (
	"context"

	"github.com/m21power/ecomm/ecomm-api/storer"
)

type Server struct {
	storer *storer.MySQLStorer
}

func NewServer(storer *storer.MySQLStorer) *Server {
	return &Server{storer: storer}
}
func (s *Server) CreateProduct(ctx context.Context, product *storer.Product) (*storer.Product, error) {
	pr, err := s.storer.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (s *Server) GetProduct(ctx context.Context, id int64) (*storer.Product, error) {
	pr, err := s.storer.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	return pr, nil

}

func (s *Server) ListProducts(ctx context.Context) ([]storer.Product, error) {
	pr, err := s.storer.ListProducts(ctx)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
func (s *Server) UpdateProduct(ctx context.Context, product *storer.Product) (*storer.Product, error) {
	return s.storer.UpdateProduct(ctx, product)
}
func (s *Server) DeleteProduct(ctx context.Context, id int64) error {
	return s.storer.DeleteProduct(ctx, id)
}
