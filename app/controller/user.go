package controller

import (
	"database/sql"
	"fmt"
	"github.com/koloo91/loginservice/app/security"
	"github.com/koloo91/loginservice/app/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

func profile() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.(security.UserContext).GetUser()
		return ctx.JSON(http.StatusOK, user)
	}
}

func getUserById(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userId := ctx.Param("id")

		foundUser, err := service.GetUserById(ctx.Request().Context(), db, userId)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user with id '%s' not found", userId))
		} else if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return ctx.JSON(http.StatusOK, foundUser)
	}
}
