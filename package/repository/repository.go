package repository

import (
	"REST_JWT"
	"github.com/jmoiron/sqlx"
)

type Auth interface {
	CreateUser(user REST_JWT.User) (int, error)
	GetUser(username string, password string) (REST_JWT.User, error)
	GetEmail(username string) (string, error)
	SetRefreshToken(refreshToken string, username string) (string, error)
	GetRefreshToken(username string) (string, error)
}

type ImportantInformation interface {
}

type Repository struct {
	Auth
	ImportantInformation
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Auth: NewAuthPostgres(db),
	}
}
