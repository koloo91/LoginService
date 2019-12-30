package componenttest

import (
	"context"
	"database/sql"
	"github.com/koloo91/loginservice/app"
	"github.com/koloo91/loginservice/app/controller"
	"github.com/koloo91/loginservice/app/model"
	"github.com/koloo91/loginservice/app/service"
	"github.com/labstack/echo/v4"
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
