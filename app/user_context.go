package app

import (
	"bitbucket.org/Koloo/lgn/app/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type UserContext struct {
	echo.Context
}

func (ctx *UserContext) GetUser() model.UserClaim {
	return ctx.Get("user").(*jwt.Token).Claims.(model.UserClaim)
}

func UserContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(&UserContext{Context: c})
	}
}
