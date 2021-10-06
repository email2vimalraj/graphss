package http

import "github.com/email2vimalraj/graphss/pkg/config"

type Server struct {
	Cfg *config.Cfg
}

func New(cfg *config.Cfg) (*Server, error) {
	return &Server{
		Cfg: cfg,
	}, nil
}
