package integration_tests

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/koloo91/loginservice/app/controller"
	"github.com/koloo91/loginservice/app/model"
	"github.com/koloo91/loginservice/app/service"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

const (
	connectionString = "postgres://postgres:@localhost/postgres?sslmode=disable"
	jwtSecret        = "s3cr3t"
)

type ComponentTestSuite struct {
	suite.Suite
	db     *sql.DB
	router *gin.Engine
}

func (suite *ComponentTestSuite) SetupSuite() {
	log.Println("Setup suite")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	driver, _ := postgres.WithInstance(db, &postgres.Config{})
	migrator, _ := migrate.NewWithDatabaseInstance("file://../migrations", "postgres", driver)
	if err := migrator.Up(); err != nil {
		log.Println(err.Error())
	}

	suite.db = db
	suite.router = controller.SetupRoutes(db, []byte(jwtSecret))
}

func (suite *ComponentTestSuite) SetupTest() {
	deleteUserStatement, _ := suite.db.Prepare("DELETE FROM app_user;")
	_, _ = deleteUserStatement.Exec()
}

func TestComponentTestSuite(t *testing.T) {
	suite.Run(t, new(ComponentTestSuite))
}

func (suite *ComponentTestSuite) createUser(name, password string) *model.UserVo {
	user, err := service.Register(context.Background(), suite.db, &model.RegisterVo{Name: name, Password: password})
	if err != nil {
		log.Fatal("error creating user", err)
	}
	return user
}
