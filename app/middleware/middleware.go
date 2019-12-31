package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/koloo91/loginservice/app/model"
)

type UserContext struct {
	*gin.Context
}

func (ctx UserContext) GetUser() model.UserClaim {
	value, _ := ctx.Get("user")
	mapClaim := value.(*jwt.Token).Claims.(jwt.MapClaims)

	return model.UserClaim{
		Id:   mapClaim["id"].(string),
		Name: mapClaim["name"].(string),
	}
}
