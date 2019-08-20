package componenttest

import (
	"bitbucket.org/Koloo/lgn/app"
	"bitbucket.org/Koloo/lgn/app/controller"
	"bitbucket.org/Koloo/lgn/app/model"
	"bitbucket.org/Koloo/lgn/app/service"
	"context"
	"database/sql"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

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
	suite.db = app.ConnectToDatabase("lgn_dba", "lgn_dba", "lgn", "lgn", "localhost", "lgn_service", "file://../migrations")

	suite.router = controller.SetupRoutes(suite.db, []byte("s3cr3t"))
}

func (suite *ComponentTestSuite) SetupTest() {
	deleteUserStatement, _ := suite.db.Prepare("DELETE FROM app_user;")
	_, _ = deleteUserStatement.Exec()
}

func (suite *ComponentTestSuite) createUser(name, password string) *model.UserVo {
	user, err := service.Register(context.Background(), suite.db, &model.RegisterVo{Name: name, Password: password})
	if err != nil {
		log.Fatal("error creating user", err)
	}
	return user
}
