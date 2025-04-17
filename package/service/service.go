package service

import (
	"REST_JWT"
	"REST_JWT/package/repository"
)

type Auth interface {
	CreateUser(user REST_JWT.User) (id int, err error)
	GenerateToken(uername string, password string, ip string) (token string, err error)
	ParseToken(token string) (username string, ip string, err error)
	GenerateRefreshToken(accessToken string) (string, error)
	ParseRefreshToken(token string) (username string, err error)
	GenerateNewToken(username string, ip string) (token string, err error)
	SendEmail(ip string, username string) (err error)
	SaveRefreshToken(refreshToken string, username string) (err error)
	VerifyToken(refreshToken string, username string) (isTrueToken bool, err error)
}

type ImportantInformation interface {
}

type Service struct {
	Auth
	ImportantInformation
}

func NewService(rps *repository.Repository) *Service {
	return &Service{
		Auth: NewAuthService(rps.Auth),
	}
}
