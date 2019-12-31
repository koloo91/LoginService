package integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func (suite *ComponentTestSuite) TestRegisterUserSuccessful() {
	body := []byte(`{"name": "kolo", "password": "Pass00"}`)

	request, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()
	suite.router.ServeHTTP(recorder, request)

	suite.Equal(http.StatusCreated, recorder.Code)

	var registerResponse map[string]interface{}
	_ = json.Unmarshal(recorder.Body.Bytes(), &registerResponse)

	suite.True(len(registerResponse["id"].(string)) > 0)
	suite.Equal("kolo", registerResponse["name"])
	suite.NotNil(registerResponse["created"])
	suite.NotNil(registerResponse["updated"])
}

func (suite *ComponentTestSuite) TestRegisterUserWithExistingShouldFail() {
	firstBody := bytes.NewBuffer([]byte(`{"name": "kolo", "password": "Pass00"}`))

	firstRecorder := httptest.NewRecorder()
	firstRequest, _ := http.NewRequest(http.MethodPost, "/api/register", firstBody)

	suite.router.ServeHTTP(firstRecorder, firstRequest)
	suite.Equal(http.StatusCreated, firstRecorder.Code)

	secondBody := bytes.NewBuffer([]byte(`{"name": "Kolo", "password": "Pass00"}`))

	secondRecorder := httptest.NewRecorder()
	secondRequest, _ := http.NewRequest("POST", "/api/register", secondBody)

	suite.router.ServeHTTP(secondRecorder, secondRequest)
	suite.Equal(http.StatusBadRequest, secondRecorder.Code)

	var registerResponse map[string]interface{}
	_ = json.Unmarshal(secondRecorder.Body.Bytes(), &registerResponse)

	suite.Equal("duplicate key value violates unique constraint \"user_name_uindex\"", registerResponse["message"])
}
