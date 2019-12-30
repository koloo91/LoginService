package controller

import (
	"database/sql"
	"fmt"
	appMiddleware "github.com/koloo91/loginservice/app/middleware"
	"github.com/koloo91/loginservice/app/model"
	"github.com/koloo91/loginservice/app/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lib/pq"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"net/http"

	_ "github.com/koloo91/loginservice/docs"
)

func SetupRoutes(db *sql.DB, jwtKey []byte) *echo.Echo {
	router := echo.New()

	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out},"service":"lgn"}` + "\n",
	}))

	router.Use(middleware.Recover())
	router.Use(middleware.CORS())

	corsConfig := middleware.DefaultCORSConfig
	corsConfig.AllowMethods = append(middleware.DefaultCORSConfig.AllowMethods, http.MethodOptions)

	router.Use(middleware.CORSWithConfig(corsConfig))

	log.Println("Setting up routes")

	router.GET("/swagger/*", echoSwagger.WrapHandler)

	{
		apiGroup := router.Group("/api")
		apiGroup.POST("/register", register(db))
		apiGroup.POST("/login", login(db, jwtKey))
		apiGroup.GET("/users/:id", getUserById(db))

		apiGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: jwtKey}))
		apiGroup.Use(appMiddleware.UserContextMiddleware)
		apiGroup.GET("/profile", profile())

		apiGroup.GET("/alive", alive())
	}

	return router
}

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

func profile() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.(appMiddleware.UserContext).GetUser()
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

func alive() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.HTML(http.StatusNoContent, "")
	}
}
