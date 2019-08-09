package componenttest

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"lgn/internal/controller"
	"log"
	"testing"
)

type ComponentTestSuite struct {
	suite.Suite
	db     *sql.DB
	router *echo.Echo
}

func TestComponentTestSuite(t *testing.T) {
	suite.Run(t, new(ComponentTestSuite))
}

func (suite *ComponentTestSuite) SetupSuite() {
	connectionString := fmt.Sprintf("host=localhost user=lgn password=lgn dbname=lgn_service sslmode=disable")
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	suite.db = db

	suite.router = controller.SetupRoutes(suite.db, []byte("s3cr3t"))
}

func (suite *ComponentTestSuite) SetupTest() {
	deleteUserStatement, _ := suite.db.Prepare("DELETE FROM app_user;")
	_, _ = deleteUserStatement.Exec()
}