package repositories

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepoInterface interface {
	GetAllProducts()
	GetDetailProduct()
}

type RepoProduct struct {
	DB *pgxpool.Pool
}

func NewProduct(db *pgxpool.Pool) *RepoProduct {
	return &RepoProduct{DB: db}
}

func (r *RepoProduct) GetAllProducts() {

}

func (r *RepoProduct) GetDetailProduct() {

}
