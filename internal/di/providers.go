package di

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/truongle2004/campus_marketplace/docs"
	"github.com/truongle2004/campus_marketplace/internal/database"
	"github.com/truongle2004/campus_marketplace/internal/modules/campus"
	"github.com/truongle2004/campus_marketplace/internal/modules/category"
	"github.com/truongle2004/campus_marketplace/internal/modules/user"
	"github.com/truongle2004/campus_marketplace/pkg/auth"
	"github.com/truongle2004/campus_marketplace/pkg/logger"
)

func ProvideLogger() (*logger.Logger, error) {
	return logger.New(logger.DefaultConfig())
}

func ProvidePort() int {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		return 8080
	}
	return port
}

func NewHTTPHandler(
	campusRoutes *campus.Routes,
	categoryRoutes *category.Routes,
	userRoutes *user.Routes,
) http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	campusRoutes.RegisterHealth(r)

	api := r.Group("/api/v1")
	campusRoutes.RegisterPublic(api)
	categoryRoutes.RegisterPublic(api)
	userRoutes.RegisterProtected(api)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func NewHTTPServer(port int, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

var AppSet = wire.NewSet(
	database.NewDatabase,
	ProvideLogger,
	ProvidePort,
	auth.Set,
	campus.Set,
	category.Set,
	user.Set,
	NewHTTPHandler,
	NewHTTPServer,
)
