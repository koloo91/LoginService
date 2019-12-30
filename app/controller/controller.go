package controller

import (
	"database/sql"
	"github.com/koloo91/loginservice/app/log"
	"github.com/koloo91/loginservice/app/security"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"

	_ "github.com/koloo91/loginservice/docs"
)

// @title Lgn Api

// @host meetya.io
// @BasePath /
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

	log.Info("Setting up routes")

	router.GET("/swagger/*", echoSwagger.WrapHandler)

	{
		internalGroup := router.Group("/internal")
		internalGroup.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
		internalGroup.GET("/alive", alive())
	}

	{
		apiGroup := router.Group("/api")
		apiGroup.POST("/register", register(db))
		apiGroup.POST("/login", login(db, jwtKey))
		apiGroup.GET("/users/:id", getUserById(db))

		apiGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: jwtKey}))
		apiGroup.Use(security.UserContextMiddleware)
		apiGroup.GET("/profile", profile())
	}

	return router
}

func alive() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.HTML(http.StatusNoContent, "")
	}
}
