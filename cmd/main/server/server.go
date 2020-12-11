package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/avian-digital-forensics/timeline-investigator/configs"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/services"

	"github.com/gorilla/handlers"
	"github.com/pacedotdev/oto/otohttp"
)

// Server holds all the dependencies for the API
type Server struct {
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

	// Register the services
	api.RegisterService(srv.router, &services.Service{})

	return nil
}

// Run the server
func (srv *Server) Run(cfg *configs.MainAPI) error {
	// Create a cors-handler
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"HEAD", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{
			"Accept",
			"Accept-Language",
			"Authorization",
			"Content-Language",
			"Content-Type",
			"Origin",
			"X-Requested-With",
		}),
	)

	// Create the http-server
	httpServer := &http.Server{
		Handler:      corsHandler(srv.router),
		Addr:         cfg.Network.IP + ":" + cfg.Network.Port,
		WriteTimeout: time.Duration(cfg.Network.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Network.ReadTimeout) * time.Second,
	}

	// Start the http-server
	log.Println("HTTP server @", httpServer.Addr)
	return httpServer.ListenAndServe()
}

// Stop the server
func (srv *Server) Stop() error {
	return nil
}
