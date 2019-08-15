package controller

import (
	"bitbucket.org/Koloo/lgn/app/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
)

func Profile() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token := ctx.Get("user").(*jwt.Token)
		return ctx.JSON(http.StatusOK, token.Claims.(model.UserClaim))
	}
}
