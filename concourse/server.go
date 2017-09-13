package summary

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Server struct
type Server struct {
	Config *Config
}

// CreateServer - creates a server
func CreateServer(config *Config) *Server {
	return &Server{Config: config}
}

// Start - starts the web server
func (s *Server) Start() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", s.Config.Index)
	router.HandleFunc("/host/{host}", s.Config.HostSummary)
	router.HandleFunc("/group/{group}", s.Config.GroupSummary)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets/")))

	return router
}
