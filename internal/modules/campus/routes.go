package campus

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) *Routes {
	return &Routes{handler: handler}
}

func (rt *Routes) RegisterPublic(r *gin.RouterGroup) {
	r.GET("/campuses", rt.handler.ListActive)
	r.GET("/campuses/:id", rt.handler.GetByID)
}

func (rt *Routes) RegisterHealth(r *gin.Engine) {
	r.GET("/health", rt.handler.Health)
}
