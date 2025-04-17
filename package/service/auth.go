package service

import (
	"REST_JWT"
	"REST_JWT/package/repository"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	tokenPass        = "SomeStrongPassword"
	refreshTokenPass = "aqsdadqdasdawd"
	salt             = "qasdqfpsdplgfdfkgjasdlk"
)

type AuthService struct {
	repo repository.Auth
}

type TokenClaim struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Ip       string `json:"ip"`
}

type TokenRefresh struct {
	jwt.StandardClaims
	AccessToken string `json:"access_token"`
}

func NewAuthService(repo repository.Auth) *AuthService {
	return &AuthService{repo: repo}
}

func (a *AuthService) CreateUser(user REST_JWT.User) (int, error) {
	user.Password = a.generatePasswordHash(user.Password)
	return a.repo.CreateUser(user)
}

func (s *AuthService) GenerateRefreshToken(accessToken string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenRefresh{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 120).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		accessToken,
	})
	return token.SignedString([]byte(refreshTokenPass))
}

func (s *AuthService) GenerateToken(username string, password string, ip string) (string, error) {
	user, err := s.repo.GetUser(username, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &TokenClaim{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Username,
		ip,
	})

	return token.SignedString([]byte(tokenPass))
}
func (s *AuthService) GenerateNewToken(username string, ip string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &TokenClaim{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		username,
		ip,
	})
	return token.SignedString([]byte(tokenPass))
}

func (a *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (a *AuthService) ParseToken(accesstoken string) (string, string, error) {
	token, err := jwt.ParseWithClaims(accesstoken, &TokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(tokenPass), nil
	})
	if err != nil {
		return "", "", err
	}

	claim, ok := token.Claims.(*TokenClaim)
	if !ok {
		return "", "", fmt.Errorf("Invalid token claim")
	}

	return claim.Username, claim.Ip, nil
}

func (a *AuthService) ParseRefreshToken(refreshtoken string) (string, error) {
	token, err := jwt.ParseWithClaims(refreshtoken, &TokenRefresh{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(refreshTokenPass), nil
	})

	if err != nil {
		return "", err
	}
	claim, ok := token.Claims.(*TokenRefresh)
	if !ok {
		return "", fmt.Errorf("Invalid token claim")
	}
	return claim.AccessToken, nil

}
func (a *AuthService) VerifyToken(refreshToken string, username string) (bool, error) {
	trueToken, err := a.repo.GetRefreshToken(username)
	if err != nil {
		return false, err
	}
	if bcrypt.CompareHashAndPassword([]byte(trueToken), []byte(hashToken(refreshToken))) != nil {
		return false, nil
	}
	return true, nil
}

func (a *AuthService) SendEmail(ip string, username string) (err error) {
	var email string
	email, err = a.repo.GetEmail(username)
	if err != nil {
		return err
	}
	fmt.Printf("Сообщение об новом ip адресе (%s) отправлено на %s\n", ip, email)
	//Отправка сообщения о новом ip-адресе
	return nil
}

func (a *AuthService) SaveRefreshToken(refreshToken string, username string) (err error) {
	refreshTokenCRYPT, err := bcrypt.GenerateFromPassword([]byte(hashToken(refreshToken)), bcrypt.MinCost)
	if err != nil {
		return err
	}
	_, err = a.repo.SetRefreshToken(string(refreshTokenCRYPT), username)
	return err
}
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
