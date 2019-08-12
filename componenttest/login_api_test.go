package componenttest

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
)

func (suite *ComponentTestSuite) TestLoginUserSuccessful() {
	suite.createUser("foo", "bar")

	body := bytes.NewBuffer([]byte(`{"name": "foo", "password": "bar"}`))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/login", body)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	suite.router.ServeHTTP(recorder, request)

	suite.Equal(http.StatusOK, recorder.Code)

	var loginResponse map[string]interface{}
	suite.Nil(json.NewDecoder(recorder.Body).Decode(&loginResponse))

	suite.True(len(loginResponse["token"].(string)) > 0)
	suite.Equal("Bearer", loginResponse["type"])
}

func (suite *ComponentTestSuite) TestLoginUserShouldFail() {
	body := bytes.NewBuffer([]byte(`{"name": "foo", "password": "bar"}`))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/login", body)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	suite.router.ServeHTTP(recorder, request)

	suite.Equal(http.StatusBadRequest, recorder.Code)

	var loginResponse map[string]interface{}
	suite.Nil(json.NewDecoder(recorder.Body).Decode(&loginResponse))

	suite.Equal("Invalid credentials", loginResponse["message"])
}
