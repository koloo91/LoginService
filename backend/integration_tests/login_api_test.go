package integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func (suite *ComponentTestSuite) TestLoginUserSuccessful() {
	suite.createUser("foo", "bar")

	body := bytes.NewBuffer([]byte(`{"name": "foo", "password": "bar"}`))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/api/login", body)

	suite.router.ServeHTTP(recorder, request)

	suite.Equal(http.StatusOK, recorder.Code)

	var loginResponse map[string]interface{}
	suite.Nil(json.NewDecoder(recorder.Body).Decode(&loginResponse))

	suite.True(len(loginResponse["accessToken"].(string)) > 0)
	suite.Equal("Bearer", loginResponse["type"])
}

func (suite *ComponentTestSuite) TestLoginUserShouldFail() {
	body := bytes.NewBuffer([]byte(`{"name": "foo", "password": "bar"}`))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/api/login", body)

	suite.router.ServeHTTP(recorder, request)

	suite.Equal(http.StatusBadRequest, recorder.Code)

	var loginResponse map[string]interface{}
	suite.Nil(json.NewDecoder(recorder.Body).Decode(&loginResponse))

	suite.Equal("invalid credentials", loginResponse["message"])
}
