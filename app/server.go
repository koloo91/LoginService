package app

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/koloo91/loginservice/app/controller"
	"github.com/koloo91/loginservice/app/log"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"os"
	"time"
)

var jwtKey []byte

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	jwtKey = []byte(getEnvOrDefault("JWT_KEY", "s3cr3t"))
}

func Start() {
	log.Info("Starting application")

	dbaUser := getEnvOrDefault("DBA_USER", "lgn_dba")
	dbaPassword := getEnvOrDefault("DBA_PASSWORD", "lgn_dba")
	dbUser := getEnvOrDefault("DB_USER", "lgn")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "lgn")
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbName := getEnvOrDefault("DB_NAME", "lgn_service")

	db := ConnectToDatabase(dbaUser, dbaPassword, dbUser, dbPassword, dbHost, dbName, "file://migrations")

	router := controller.SetupRoutes(db, jwtKey)

	port := getEnvOrDefault("PORT", "8080")
	log.Infof("Starting http server on port %s", port)

	log.Fatal(router.Start(fmt.Sprintf(":%s", port)))
}

func ConnectToDatabase(dbaUser, dbaPassword, dbUser, dbPassword, host, dbName, migrationFilesPath string) *sql.DB {
	dbaConnectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, dbaUser, dbaPassword, dbaUser)
	dbaDb, err := sql.Open("postgres", dbaConnectionString)
	if err != nil {
		log.Fatalf("Error opening dba database connection '%s'", err.Error())
		return nil
	}

	if err := dbaDb.Ping(); err != nil {
		log.Fatalf("Error pinging database with dba user '%s'", err.Error())
		return nil
	}

	dbaDriver, err := postgres.WithInstance(dbaDb, &postgres.Config{})
	dbaMigrations, _ := migrate.NewWithDatabaseInstance(fmt.Sprintf("%s/dba", migrationFilesPath), "postgres", dbaDriver)
	if err := dbaMigrations.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error migrating database ", err)
	}

	_ = dbaDb.Close()

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Error opening database connection '%s'", err.Error())
		return nil
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database '%s'", err.Error())
		return nil
	}

	dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
	dbMigrations, _ := migrate.NewWithDatabaseInstance(fmt.Sprintf("%s/lgn", migrationFilesPath), "postgres", dbDriver)
	if err := dbMigrations.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error migrating database ", err)
	}

	log.Infof("Connected to database '%s' with user '%s'", host, dbUser)

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
