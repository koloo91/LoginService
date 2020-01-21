package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/koloo91/loginservice/controller"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	dbHost     = getEnvOrDefault("DB_HOST", "localhost")
	dbUser     = getEnvOrDefault("DB_USER", "postgres")
	dbPassword = getEnvOrDefault("DB_PASSWORD", "")
	dbName     = getEnvOrDefault("DB_NAME", "postgres")

	jwtKey                    = []byte(getEnvOrDefault("JWT_KEY", "s3cr3t"))
	validateExpirationDate, _ = strconv.ParseBool(getEnvOrDefault("VALIDATE_EXPIRATION_DATE", "true"))
)

func main() {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbUser, dbPassword, dbHost, 5432, dbName)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("unable to open database connection: %s", err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("unable to ping database: %s", err.Error())
	}

	defer db.Close()

	driver, _ := postgres.WithInstance(db, &postgres.Config{})
	migrator, _ := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err := migrator.Up(); err != nil {
		log.Println(err.Error())
	}

	db.SetConnMaxLifetime(1 * time.Hour)
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(2)

	router := controller.SetupRoutes(db, jwtKey, validateExpirationDate)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.File("./assets/index.html")
	})

	router.Static("/app", "./assets")

	log.Fatal(router.Run())
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
