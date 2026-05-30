package category

import "github.com/google/wire"

var Set = wire.NewSet(
	NewRepo,
	NewService,
	NewHandler,
	NewRoutes,
)
