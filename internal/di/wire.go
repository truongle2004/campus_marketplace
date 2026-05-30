//go:build wireinject
// +build wireinject

package di

import (
	"net/http"

	"github.com/google/wire"
	"github.com/truongle2004/campus_marketplace/pkg/logger"
)

type App struct {
	Server *http.Server
	Logger *logger.Logger
}

func InitializeApp() (*App, error) {
	wire.Build(
		AppSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil
}
