package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

func (suite *ComponentTestSuite) TestShouldReturnUserProfileSuccessful() {
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

	profileRecorder := httptest.NewRecorder()
	profileRequest, _ := http.NewRequest(http.MethodGet, "/api/profile", nil)
	profileRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", loginResponse["accessToken"]))
	suite.router.ServeHTTP(profileRecorder, profileRequest)

	var profileResponse map[string]interface{}
	suite.Nil(json.NewDecoder(profileRecorder.Body).Decode(&profileResponse))

	suite.True(len(profileResponse["id"].(string)) > 0)
	suite.Equal("foo", profileResponse["name"].(string))
}
