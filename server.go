package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	router    *mux.Router
	scheduler *scheduler
	state     *State
}

func newServer(listenAddr string, scheduler *scheduler, state *State) *server {
	return &server{
		router:    mux.NewRouter(),
		scheduler: scheduler,
		state:     state,
	}
}

func (server *server) start(bindAddr string) error {

	server.router.HandleFunc("/api/state", server.state.httpHandlerState).Methods("GET")
	server.router.HandleFunc("/api/points", server.state.httpHandlerPoints).Methods("GET")
	server.router.HandleFunc("/api/tasks", server.scheduler.httpHandlerTasks).Methods("GET")
	server.router.HandleFunc("/api/task/{id}", server.scheduler.httpHandlerTask).Methods("GET", "DELETE")

	return http.ListenAndServe(bindAddr, server.router)
}
