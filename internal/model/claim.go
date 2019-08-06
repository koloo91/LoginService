package model

import "github.com/dgrijalva/jwt-go"

type Claim struct {
	jwt.StandardClaims
	Id   string `json:"id"`
	Name string `json:"name"`
}
