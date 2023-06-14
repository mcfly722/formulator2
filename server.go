package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	router    *mux.Router
	scheduler *scheduler
}

func NewServer(listenAddr string, scheduler *scheduler) *server {
	return &server{
		router:    mux.NewRouter(),
		scheduler: scheduler,
	}

}

func (server *server) start(bindAddr string) error {

	server.router.HandleFunc("/api/points", server.scheduler.state.HTTPHandlerPoints)
	server.router.HandleFunc("/api/tasks", server.scheduler.HTTPHandlerTasks)
	server.router.HandleFunc("/api/task/new", server.scheduler.HTTPHandlerNewTask)
	return http.ListenAndServe(bindAddr, server.router)
}
