package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/koloo91/loginservice/app/model"
	"github.com/labstack/echo/v4"
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
