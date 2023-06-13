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
	server.router.HandleFunc("/ListPoints", server.scheduler.state.HTTPHandlerListPoints)
	server.router.HandleFunc("/Scheduler/ListTasks", server.scheduler.HTTPHandlerListTasks)
	server.router.HandleFunc("/Scheduler/NewTask", server.scheduler.HTTPHandlerNewTask)
	return http.ListenAndServe(bindAddr, server.router)
}
