//go:build wireinject
// +build wireinject

package http

import (
	"github.com/email2vimalraj/graphss/pkg/config"
	"github.com/google/wire"
)

var providerSet = wire.NewSet(config.NewCfg, New)

func InitializeServer(configFilePath string) (*Server, error) {
	wire.Build(providerSet)
	return &Server{}, nil
}
