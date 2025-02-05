package common

import "github.com/golang-jwt/jwt/v5"

type CONTEXT_KEY string

const (
	TOKEN_EXPIRY_MINUTES = 30
	USER_ID_KEY CONTEXT_KEY = "user_id"
)

type User struct{
	UserData
	ID 			uint64		`json:"id"`
	CreatedAt	string		`json:"created_at"`
}

type UserData struct{
	Name		string		`json:"name"`
	Email		string		`json:"email"`
	Password 	string		`json:"password"`
}

type Claims struct{
	jwt.RegisteredClaims
	UserID		uint64		`json:"user_id"`
}