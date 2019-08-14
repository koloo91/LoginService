package controller

import (
	"bitbucket.org/Koloo/lgn/app/model"
	"bitbucket.org/Koloo/lgn/app/service"
	"database/sql"
	"github.com/labstack/echo"
	"github.com/lib/pq"
	"net/http"
)

func Register(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var registerVo model.RegisterVo
		if err := ctx.Bind(&registerVo); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid json")
		}

		userVo, err := service.Register(ctx.Request().Context(), db, &registerVo)
		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				return echo.NewHTTPError(http.StatusBadRequest, err.Message)
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return ctx.JSON(http.StatusCreated, userVo)
	}
}

func Login(db *sql.DB, jwtKey []byte) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var loginVo model.LoginVo
		if err := ctx.Bind(&loginVo); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid json")
		}

		token, err := service.Login(ctx.Request().Context(), db, jwtKey, &loginVo)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid credentials")
		}

		return ctx.JSON(http.StatusOK, map[string]string{"token": token, "type": "Bearer"})
	}
}
