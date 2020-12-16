package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/avian-digital-forensics/timeline-investigator/configs"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/services"

	"github.com/pacedotdev/oto/otohttp"
)

// Server holds all the dependencies for the API
type Server struct {
	http   *http.Server
	router *otohttp.Server
	ctx    context.Context
}

// New creates a new server
func New(ctx context.Context) *Server {
	return &Server{
		ctx:    ctx,
		router: otohttp.NewServer(),
	}
}

// Initialize the server
func (srv *Server) Initialize(cfg *configs.MainAPI) error {
	// Set the base-path for the oto-server
	srv.router.Basepath = "/api/"
	http.Handle("/api/", srv.router)

	// endpoint for for gke-healthchecks
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	// Register the services
	api.RegisterService(srv.router, &services.Service{})

	return nil
}

// Run the server
func (srv *Server) Run(cfg *configs.MainAPI) error {
	// Create the http-server
	srv.http = &http.Server{
		Addr:         cfg.Network.IP + ":" + cfg.Network.Port,
		WriteTimeout: time.Duration(cfg.Network.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Network.ReadTimeout) * time.Second,
	}

	// Start the http-server
	log.Println("HTTP server @", srv.http.Addr)
	return srv.http.ListenAndServe()
}

// Stop the server
func (srv *Server) Stop() error {
	ctx, cancel := context.WithTimeout(srv.ctx, 2*time.Second)
	defer cancel()

	return srv.http.Shutdown(ctx)
}
