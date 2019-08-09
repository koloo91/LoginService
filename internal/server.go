package internal

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"lgn/internal/controller"
	"lgn/internal/logging"
	"os"
	"time"
)

var jwtKey []byte

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	jwtKey = []byte(getEnvOrDefault("JWT_KEY", "s3cr3t"))
}

func Start() {
	logging.AppLogger.Info("Starting application")

	db := connectToDatabase()

	router := controller.SetupRoutes(db, jwtKey)

	port := getEnvOrDefault("PORT", "8080")
	logging.AppLogger.Infof("Starting http server on port %s", port)

	logging.AppLogger.Fatal(router.Start(fmt.Sprintf(":%s", port)))
}

func connectToDatabase() *sql.DB {
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbUser := getEnvOrDefault("DB_USER", "lgn")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "lgn")
	dbName := getEnvOrDefault("DB_NAME", "lgn_service")

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName)

	logrus.Infof("Using connection string '%s'", connectionString)

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		logging.AppLogger.Fatalf("Error opening database connection '%s'", err.Error())
		return nil
	}

	if err := db.Ping(); err != nil {
		logging.AppLogger.Fatalf("Error pinging database '%s'", err.Error())
		return nil
	}

	logging.AppLogger.Infof("Connected to database '%s' with user '%s'", dbHost, dbUser)

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
