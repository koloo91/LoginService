package componenttest

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo"
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

	var registerResponse map[string]interface{}
	suite.Nil(json.NewDecoder(recorder.Body).Decode(&registerResponse))

	suite.True(len(registerResponse["id"].(string)) > 0)
	suite.Equal("kolo", registerResponse["name"])
	suite.NotNil(registerResponse["created"])
	suite.NotNil(registerResponse["updated"])
}

func (suite *ComponentTestSuite) TestRegisterUserWithExistingShouldFail() {
	firstBody := bytes.NewBuffer([]byte(`{"name": "kolo", "password": "Pass00"}`))

	firstRecorder := httptest.NewRecorder()
	firstRequest, _ := http.NewRequest("POST", "/api/register", firstBody)
	firstRequest.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	suite.router.ServeHTTP(firstRecorder, firstRequest)
	suite.Equal(http.StatusCreated, firstRecorder.Code)

	secondBody := bytes.NewBuffer([]byte(`{"name": "Kolo", "password": "Pass00"}`))

	secondRecorder := httptest.NewRecorder()
	secondRequest, _ := http.NewRequest("POST", "/api/register", secondBody)
	secondRequest.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	suite.router.ServeHTTP(secondRecorder, secondRequest)
	suite.Equal(http.StatusBadRequest, secondRecorder.Code)

	var registerResponse map[string]interface{}
	suite.Nil(json.NewDecoder(secondRecorder.Body).Decode(&registerResponse))

	suite.Equal("duplicate key value violates unique constraint \"user_name_uindex\"", registerResponse["message"])
}
