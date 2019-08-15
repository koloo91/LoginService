package model

import "github.com/dgrijalva/jwt-go"

type UserClaim struct {
	jwt.StandardClaims
	Id   string `json:"id"`
	Name string `json:"name"`
}
