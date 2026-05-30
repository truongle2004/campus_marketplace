package user

import (
	"github.com/gin-gonic/gin"
	"github.com/truongle2004/campus_marketplace/pkg/auth"
)

type Routes struct {
	handler    *Handler
	middleware *auth.Middleware
}

func NewRoutes(handler *Handler, middleware *auth.Middleware) *Routes {
	return &Routes{handler: handler, middleware: middleware}
}

func (rt *Routes) RegisterProtected(r *gin.RouterGroup) {
	me := r.Group("/users/me", rt.middleware.RequireAuth())
	{
		me.POST("/bootstrap", rt.handler.Bootstrap)
		me.GET("", rt.handler.GetMe)
		me.PATCH("", rt.handler.UpdateMe)
		me.DELETE("", rt.handler.DeactivateMe)
	}
}
