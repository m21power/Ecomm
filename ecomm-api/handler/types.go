package handler

import "time"

type ProductReq struct {
	Name         string  `json:"name"`
	Image        string  `json:"image"`
	Category     string  `json:"category"`
	Description  string  `json:"description"`
	Rating       int     `json:"rating"`
	NumReviews   int     `json:"num_reviews"`
	Price        float64 `json:"price"`
	CountInStock int64   `json:"count_in_stock"`
}
type ProductRes struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Image        string     `json:"image"`
	Category     string     `json:"category"`
	Description  string     `json:"description"`
	Rating       int        `json:"rating"`
	NumReviews   int        `json:"num_reviews"`
	Price        float64    `json:"price"`
	CountInStock int64      `json:"count_in_stock"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
