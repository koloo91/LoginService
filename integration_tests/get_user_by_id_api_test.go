package integration_tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

func (suite *ComponentTestSuite) TestGetUserByIdSuccessful() {
	user := suite.createUser("foo", "bar")

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/users/%s", user.Id), nil)

	suite.router.ServeHTTP(recorder, request)

	suite.Equal(http.StatusOK, recorder.Code)

	var byIdResponse map[string]interface{}
	suite.Nil(json.NewDecoder(recorder.Body).Decode(&byIdResponse))

	suite.Equal(user.Id, byIdResponse["id"])
	suite.Equal("foo", byIdResponse["name"])
}

func (suite *ComponentTestSuite) TestGetUserByIdShouldReturnNotFound() {

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/api/users/foo", nil)

	suite.router.ServeHTTP(recorder, request)

	suite.Equal(http.StatusNotFound, recorder.Code)

	var byIdResponse map[string]interface{}
	suite.Nil(json.NewDecoder(recorder.Body).Decode(&byIdResponse))

	suite.Equal("user with id 'foo' not found", byIdResponse["message"])
}
