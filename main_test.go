package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerHandleRoot(t *testing.T) {
	server := &Server{}

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	recorder := httptest.NewRecorder()
	server.handleRoot(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(
		t,
		"Hello world from Go net http <3",
		recorder.Body.String(),
	)
}

func TestServerCreateUser(t *testing.T) {
	server := &Server{
		userCache:  make(map[int]User),
		nextUserId: 1,
	}
	body := strings.NewReader(`{"name":"Yash"}`)
	request := httptest.NewRequest(
		http.MethodPost,
		"/users",
		body,
	)
	recorder := httptest.NewRecorder()
	server.createUser(recorder, request)
	assert.Equal(t, http.StatusNoContent, recorder.Code)
	assert.Equal(
		t,
		User{
			Id:   1,
			Name: "Yash",
		},
		server.userCache[1],
	)
	// TODO: Test multiple cases
}

func TestServer_getUserById(t *testing.T) {

	server := &Server{
		userCache: map[int]User{
			1: {
				Id:   1,
				Name: "Yash",
			},
		},
		nextUserId: 2,
	}
	request := httptest.NewRequest(
		http.MethodGet,
		"/users/1",
		nil,
	)
	request.SetPathValue("id", "1")
	recorder := httptest.NewRecorder()
	server.getUserById(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(
		t,
		`{"id":1,"name":"Yash"}`,
		recorder.Body.String(),
	)
}
