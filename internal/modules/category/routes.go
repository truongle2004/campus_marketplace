package category

import "github.com/gin-gonic/gin"

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) *Routes {
	return &Routes{handler: handler}
}

func (rt *Routes) RegisterPublic(r *gin.RouterGroup) {
	r.GET("/categories", rt.handler.ListTopLevel)
	r.GET("/categories/:id/children", rt.handler.ListChildren)
}
