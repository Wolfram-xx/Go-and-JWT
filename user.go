package REST_JWT

type User struct {
	Username     string `json:"username" binding:"required" db:"username"`
	RefreshToken string `json:"refresh_token" db:"refresh_token"`
	Password     string `json:"password" binding:"required" db:"password_hash"`
	Id           int    `json:"-" db:"id"`
	Email        string `json:"email" db:"email"`
}
