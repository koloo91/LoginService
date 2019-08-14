package controller

import (
	"database/sql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
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

	logrus.Info("Setting up routes")

	{
		internalGroup := router.Group("/internal")
		internalGroup.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	}

	{
		apiGroup := router.Group("/api")
		apiGroup.POST("/register", Register(db))
		apiGroup.POST("/login", Login(db, jwtKey))

		apiGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: jwtKey}))
		apiGroup.GET("/profile", Profile())
	}

	return router
}
