package controller

import (
	"database/sql"
	"github.com/koloo91/loginservice/app/model"
	"github.com/koloo91/loginservice/app/service"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"net/http"
)

// Register godoc
// @Summary Registers a new user
// @ID register
// @Accept json
// @Produce json
// @Param registerVo body model.RegisterVo true "register json"
// @Success 200 {object} model.UserVo
// @Failure 400 {object} model.HttpError
// @Router /api/register [post]
func register(db *sql.DB) echo.HandlerFunc {
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

// Login godoc
// @Summary Login a user
// @ID login
// @Accept json
// @Produce json
// @Param loginVo body model.LoginVo true "login json"
// @Success 200 {object} model.LoginResultVo
// @Failure 400 {object} model.HttpError
// @Router /api/login [post]
func login(db *sql.DB, jwtKey []byte) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var loginVo model.LoginVo
		if err := ctx.Bind(&loginVo); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid json")
		}

		token, err := service.Login(ctx.Request().Context(), db, jwtKey, &loginVo)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid credentials")
		}

		return ctx.JSON(http.StatusOK, model.LoginResultVo{Token: token, Type: "Bearer"})
	}
}
