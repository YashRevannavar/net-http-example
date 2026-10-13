package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var cacheMutex sync.RWMutex

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var userCache = make(map[int]User)
var nextUserId = 1

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleRoot)
	mux.HandleFunc("POST /users", createUser)
	mux.HandleFunc("GET /users/{id}", getUserById)
	mux.HandleFunc("DELETE /users/{id}", deleteUserById)

	fmt.Println("Go net http Server started")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleRoot(
	w http.ResponseWriter,
	r *http.Request,
) {
	fmt.Fprintf(w, "Hello world from Go net http <3")
}

func createUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}
	if user.Name == "" {
		http.Error(
			w,
			"Name is required",
			http.StatusBadRequest,
		)
		return
	}
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	user.Id = nextUserId
	userCache[user.Id] = user
	nextUserId++

	w.WriteHeader(http.StatusNoContent)

}

func getUserById(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	resultUser, ok := userCache[id]
	if !ok {
		http.Error(
			w,
			"User not found",
			http.StatusNotFound,
		)
		return
	}

	jsonUser, err := json.Marshal(resultUser)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUser)

}

func deleteUserById(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}
	cacheMutex.Lock()
	defer cacheMutex.RUnlock()
	if _, ok := userCache[id]; !ok {
		http.Error(
			w,
			"User not found",
			http.StatusNotFound,
		)
		return
	}
	delete(userCache, id)
	w.WriteHeader(http.StatusNoContent)
}
