package http

import (
	"fmt"
	"net"
	"net/http"
	"path"
	"strings"

	"github.com/email2vimalraj/graphss/pkg/config"
	"github.com/email2vimalraj/graphss/pkg/domains"
	"github.com/gorilla/mux"
)

// Server represents an HTTP server. It is meant to wrap all the HTTP functionality
// used by the application.
type Server struct {
	ln     net.Listener
	server *http.Server
	router *mux.Router

	Cfg *config.Cfg
}

// New returns a new instance of Server.
func New(cfg *config.Cfg) (*Server, error) {
	s := &Server{
		server: &http.Server{},
		router: mux.NewRouter(),
		Cfg:    cfg,
	}

	// Our router is wrapped by another function handler to perform some
	// middleware-like tasks that cannot be performed by actual middleware.
	// This includes changing route paths for JSON endpoints & overriding methods.
	s.server.Handler = http.HandlerFunc(s.serveHTTP)

	// Setup endpoint to display deployed version.
	s.router.HandleFunc("/debug/version", s.handleVersion).Methods(http.MethodGet)

	return s, nil
}

func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	// Override method for forms passing "_method" value.
	if r.Method == http.MethodPost {
		switch v := r.PostFormValue("_method"); v {
		case http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete:
			r.Method = v
		}
	}

	// Override content-type for certain extensions.
	// This allows us to easily cURL API endpoints with a ".json" extension
	// instead of having to explicitly set Content-type & Accept headers.
	// The extensions are removed so they don't appear in the routers.
	switch ext := path.Ext(r.URL.Path); ext {
	case ".json":
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Accept", "application/json")
		r.URL.Path = strings.TrimSuffix(r.URL.Path, ext)
	}

	// Delegate remaining HTTP handling to the gorilla router.
	s.router.ServeHTTP(w, r)
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("Version: %s\nBuild: %s\n", domains.Version, domains.Commit)))
}

func (s *Server) Open() (err error) {
	serveAddr := net.JoinHostPort(s.Cfg.HTTP.Addr, s.Cfg.HTTP.Port)

	// Open a listener on our bind address
	if s.ln, err = net.Listen("tcp", serveAddr); err != nil {
		return err
	}

	fmt.Printf("Serving on %s", serveAddr)

	// Begin serving requests on the listener. We use Serve() instead of
	// ListenAndServe() because it allows us to check for listen errors (such
	// as trying to use an already open port) synchronously.
	go s.server.Serve(s.ln)

	return nil
}
