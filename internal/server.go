package internal

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"lgn/internal/controller"
	"net/http"
	"os"
	"time"
)

var AppLogger *logrus.Entry
var jwtKey []byte

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	AppLogger = logrus.WithField("service", "lgn")
	jwtKey = []byte(getEnvOrDefault("JWT_KEY", "s3cr3t"))
}

func Start() {
	AppLogger.Info("Starting application")

	db := connectToDatabase()

	e := echo.New()
	SetupRouter(e, db)

	port := getEnvOrDefault("PORT", "8080")
	AppLogger.Infof("Starting http server on port %s", port)

	AppLogger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}

func SetupRouter(e *echo.Echo, db *sql.DB) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out},"service":"lgn"}` + "\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	corsConfig := middleware.DefaultCORSConfig
	corsConfig.AllowMethods = append(middleware.DefaultCORSConfig.AllowMethods, http.MethodOptions)

	e.Use(middleware.CORSWithConfig(corsConfig))

	AppLogger.Info("Setting up routes")

	{
		internalGroup := e.Group("/internal")
		internalGroup.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	}

	{
		apiGroup := e.Group("/api")
		apiGroup.POST("/api/register", controller.Register(db))
		apiGroup.POST("/api/login", controller.Login(db, jwtKey))
	}
}

func connectToDatabase() *sql.DB {
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbUser := getEnvOrDefault("DB_USER", "lgn")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "lgn")
	dbName := getEnvOrDefault("DB_NAME", "lgn")

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName)

	logrus.Infof("Using connection string '%s'", connectionString)

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		AppLogger.Fatalf("Error opening database connection '%s'", err.Error())
		return nil
	}

	if err := db.Ping(); err != nil {
		AppLogger.Fatalf("Error pinging database '%s'", err.Error())
		return nil
	}

	AppLogger.Infof("Connected to database '%s' with user '%s'", dbHost, dbUser)

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return db
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
