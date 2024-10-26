package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
)

type APIResponse struct {
	Statuscode int         `json:"statuscode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Error      string      `json:"error"`
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

var repeatedRequest SignupRequest

func TestAuthHome(t *testing.T) {
	// Initialize the router
	r := SetupRouter()
	req, err := http.NewRequest("GET", "/auth", nil)
	if err != nil {
		t.Fatal("Failed to create request: ", err)
	}

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, req)

	message := fmt.Sprintf("Expected 200 OK got %d", recorder.Code)

	assert.Equal(t, recorder.Code, http.StatusOK, message)

	contentTypeHeader := recorder.Header().Get("Content-Type")
	expectedContentTypeHeader := "application/json"

	message = fmt.Sprintf("Expected to receive json data got %v", contentTypeHeader)

	assert.Equal(t, contentTypeHeader, expectedContentTypeHeader, message)

	var response APIResponse

	err = json.Unmarshal(recorder.Body.Bytes(), &response)

	if err != nil {
		t.Fatal("Failed to unmarshal response: ", err)
	}

	expectedMessage := "You've hit api route"
	message = fmt.Sprintf("Expected to receive '%s' got %s", expectedMessage, response.Message)

	assert.Equal(t, expectedMessage, response.Message, message)
}

func TestAuthSignupGood(t *testing.T) {
	// Initialize the router
	r := SetupRouter()

	username := faker.Username()
	password := faker.Password()
	email := faker.Email()

	repeatedRequest.Email = email
	repeatedRequest.Username = username
	repeatedRequest.Password = password

	signupJson, err := json.Marshal(SignupRequest{
		Username: username,
		Password: password,
		Email:    email,
	})

	if err != nil {
		t.Fatal("Cannot marshal request json")
	}

	req, err := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(signupJson))

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, req)

	expectedStatusCode := http.StatusOK
	message := fmt.Sprintf("Expected 200 OK got %d", recorder.Code)
	recorderStatusCode := recorder.Code

	assert.Equal(t, expectedStatusCode, recorderStatusCode, message)

	var response APIResponse

	err = json.Unmarshal(recorder.Body.Bytes(), &response)

	if err != nil {
		t.Fatal("Failed to unmarshal response", err)
	}

	expectedMessage := "User Added"
	gotMessage := response.Message
	message = fmt.Sprintf("Expected '%s' got '%s'", expectedMessage, gotMessage)

	assert.Equal(t, expectedMessage, gotMessage, message)
}

func TestAuthSignupRepeatedDetails(t *testing.T) {

	r := SetupRouter()
	repeatedRequestJson, err := json.Marshal(repeatedRequest)
	if err != nil {
		t.Fatal("Cannot marshal request json")
	}
	req, err := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(repeatedRequestJson))

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, req)

	expectedStatusCode := http.StatusBadRequest
	gotStatusCode := recorder.Code
	message := fmt.Sprintf("Expected %d got %d", expectedStatusCode, gotStatusCode)

	assert.Equal(t, expectedStatusCode, gotStatusCode, message)
}

func TestAuthLoginGood(t *testing.T) {
	r := SetupRouter()
	loginRequest := LoginRequest{
		Username: repeatedRequest.Username,
		Email:    repeatedRequest.Email,
		Password: repeatedRequest.Password,
	}

	loginRequestJson, err := json.Marshal(loginRequest)

	if err != nil {
		t.Fatal("Cannot marshal request json")
	}

	req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginRequestJson))
	if err != nil {
		t.Fatal("Error creating request", err)
	}

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, req)

	expectedStatusCode := http.StatusOK
	gotStatusCode := recorder.Code
	message := fmt.Sprintf("Expected status code %d got %d", expectedStatusCode, gotStatusCode)

	var response APIResponse

	err = json.Unmarshal(recorder.Body.Bytes(), &response)

	if err != nil {
		t.Fatal("Cannot marshal request json")
	}

	fmt.Println(response.Message, loginRequest.Email, loginRequest.Username)

	assert.Equal(t, expectedStatusCode, gotStatusCode, message)
}
