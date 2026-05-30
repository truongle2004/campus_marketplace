package category

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/truongle2004/campus_marketplace/pkg/respond"
)

type Handler struct {
	service CategoryService
}

func NewHandler(service CategoryService) *Handler {
	return &Handler{service: service}
}

// ListTopLevel godoc
//
//	@Summary		List top-level categories
//	@Tags			categories
//	@Produce		json
//	@Success		200	{array}		Response
//	@Failure		500	{object}	respond.ErrorBody
//	@Router			/categories [get]
func (h *Handler) ListTopLevel(c *gin.Context) {
	items, err := h.service.ListTopLevel(c.Request.Context())
	if err != nil {
		respond.InternalError(c)
		return
	}
	respond.OK(c, items)
}

// ListChildren godoc
//
//	@Summary		List subcategories by parent
//	@Tags			categories
//	@Produce		json
//	@Param			id	path		string	true	"Parent category ID"
//	@Success		200	{array}		Response
//	@Failure		400	{object}	respond.ErrorBody
//	@Failure		404	{object}	respond.ErrorBody
//	@Failure		500	{object}	respond.ErrorBody
//	@Router			/categories/{id}/children [get]
func (h *Handler) ListChildren(c *gin.Context) {
	items, err := h.service.ListByParentID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if err.Error() == "invalid id" {
			respond.BadRequest(c, "invalid category id")
			return
		}
		if errors.Is(err, ErrNotFound) {
			respond.NotFound(c, "category not found")
			return
		}
		respond.InternalError(c)
		return
	}
	respond.OK(c, items)
}
