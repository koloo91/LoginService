package security

import (
	"bitbucket.org/Koloo/lgn/app/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type UserContext struct {
	echo.Context
}

func (ctx UserContext) GetUser() model.UserClaim {
	mapClaim := ctx.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)
	return model.UserClaim{
		Id:   mapClaim["id"].(string),
		Name: mapClaim["name"].(string),
	}
}

func UserContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(UserContext{Context: c})
	}
}
