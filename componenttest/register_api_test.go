package componenttest

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo"
	"lgn/internal/model"
	"net/http"
	"net/http/httptest"
)

func (suite *ComponentTestSuite) TestRegisterUserSuccessful() {
	body := bytes.NewBuffer([]byte(`{"name": "kolo", "password": "Pass00"}`))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/register", body)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	suite.router.ServeHTTP(recorder, request)

	suite.Equal(http.StatusCreated, recorder.Code)

	var userVo model.UserVo
	suite.Nil(json.NewDecoder(recorder.Body).Decode(&userVo))

	suite.True(len(userVo.Id) > 0)
	suite.Equal("kolo", userVo.Name)
	suite.NotNil(userVo.Created)
	suite.NotNil(userVo.Updated)
}
