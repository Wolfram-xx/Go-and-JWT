package repository

import (
	"REST_JWT"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (p *AuthPostgres) CreateUser(user REST_JWT.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, password_hash, email) VALUES ($1, $2, $3) RETURNING id", usersTable)
	row := p.db.QueryRow(query, user.Username, user.Password, user.Email)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (REST_JWT.User, error) {
	var user REST_JWT.User
	query := fmt.Sprintf("SELECT username FROM %s WHERE username = $1 AND password_hash = $2", usersTable)
	err := r.db.Get(&user, query, username, password)
	return user, err
}

func (r *AuthPostgres) GetEmail(username string) (string, error) {
	var user REST_JWT.User
	query := fmt.Sprintf("SELECT email FROM %s WHERE username = $1", usersTable)
	err := r.db.Get(&user, query, username)
	return user.Email, err
}

func (r *AuthPostgres) SetRefreshToken(refreshToken string, username string) (string, error) {
	var user REST_JWT.User
	query := fmt.Sprintf("UPDATE %s SET refresh_token = $1 WHERE username = $2", usersTable)
	_, err := r.db.Exec(query, refreshToken, username)
	return user.RefreshToken, err
}
func (r *AuthPostgres) GetRefreshToken(username string) (string, error) {
	var user REST_JWT.User
	query := fmt.Sprintf("SELECT refresh_token FROM %s WHERE username = $1", usersTable)
	err := r.db.Get(&user, query, username)
	return user.RefreshToken, err

}
