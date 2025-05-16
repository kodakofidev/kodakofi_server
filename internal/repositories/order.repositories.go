package repositories

import "github.com/jackc/pgx/v5/pgxpool"

type OrderRepoInterface interface {
	CreateOrder()
}

type RepoOrder struct {
	DB *pgxpool.Pool
}

func NewOrder(db *pgxpool.Pool) *RepoOrder {
	return &RepoOrder{DB: db}
}

func (r *RepoOrder) CreateOrder() {

}
