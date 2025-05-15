package repositories

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepoInterface interface {
	Register()
	Login()
	Logout()
}

type RepoAuth struct {
	DB *pgxpool.Pool
}

func NewUser(db *pgxpool.Pool) *RepoAuth {
	return &RepoAuth{DB: db}
}

func (r *RepoAuth) Register() {

}

func (r *RepoAuth) Login() {

}

func (r *RepoAuth) Logout() {

}
