package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Server struct {
	userCache  map[int]User
	mutexCache sync.RWMutex
	nextUserId int
}

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	server := &Server{
		userCache:  make(map[int]User),
		nextUserId: 1,
	}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", server.handleRoot)
	mux.HandleFunc("POST /users", server.createUser)
	mux.HandleFunc("GET /users/{id}", server.getUserById)
	mux.HandleFunc("DELETE /users/{id}", server.deleteUserById)

	fmt.Println("Go net http Server started")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func (s *Server) handleRoot(
	w http.ResponseWriter,
	r *http.Request,
) {
	fmt.Fprintf(w, "Hello world from Go net http <3")
}

func (s *Server) createUser(
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
	s.mutexCache.Lock()
	defer s.mutexCache.Unlock()

	user.Id = s.nextUserId
	s.userCache[user.Id] = user
	s.nextUserId++

	w.WriteHeader(http.StatusNoContent)

}

func (s *Server) getUserById(
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
	s.mutexCache.RLock()
	defer s.mutexCache.RUnlock()
	resultUser, ok := s.userCache[id]
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

func (s *Server) deleteUserById(
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
	s.mutexCache.Lock()
	defer s.mutexCache.Unlock()
	if _, ok := s.userCache[id]; !ok {
		http.Error(
			w,
			"User not found",
			http.StatusNotFound,
		)
		return
	}
	delete(s.userCache, id)
	w.WriteHeader(http.StatusNoContent)
}
