package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/avian-digital-forensics/timeline-investigator/configs"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/api"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/authentication"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/datastore"
	"github.com/avian-digital-forensics/timeline-investigator/pkg/filestore"
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
	auth, err := authentication.New(srv.ctx, cfg.Authentication.CredentialsFile, cfg.Authentication.APIKey)
	if err != nil {
		return err
	}

	db, err := datastore.NewService(cfg.DB.URLs...)
	if err != nil {
		return err
	}

	filestore, err := filestore.New(cfg.Filestore.BasePath)
	if err != nil {
		return err
	}

	// Set the base-path for the oto-server
	srv.router.Basepath = "/api/"
	http.Handle("/api/", srv.router)

	// endpoint for for gke-healthchecks
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	// Register the services
	caseService := services.NewCaseService(db, auth)
	api.RegisterCaseService(srv.router, caseService)
	api.RegisterEventService(srv.router, services.NewEventService(db, caseService))
	api.RegisterLinkService(srv.router, services.NewLinkService(db, caseService))
	api.RegisterFileService(srv.router, services.NewFileService(db, filestore, caseService))
	api.RegisterEntityService(srv.router, services.NewEntityService(db, caseService))
	api.RegisterProcessService(srv.router, &services.ProcessService{})

	// Only create the TestService if it is a test-run
	if cfg.Test.Run {
		api.RegisterTestService(srv.router, services.NewTestService(auth, cfg.Test.Secret))
	}

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
