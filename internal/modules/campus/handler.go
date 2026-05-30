package campus

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truongle2004/campus_marketplace/pkg/respond"
)

type Handler struct {
	service CampusService
}

func NewHandler(service CampusService) *Handler {
	return &Handler{service: service}
}

// ListActive godoc
//
//	@Summary		List active campuses
//	@Description	Returns all active campuses for onboarding and filtering.
//	@Tags			campuses
//	@Produce		json
//	@Success		200	{array}		Response
//	@Failure		500	{object}	respond.ErrorBody
//	@Router			/campuses [get]
func (h *Handler) ListActive(c *gin.Context) {
	items, err := h.service.ListActive(c.Request.Context())
	if err != nil {
		respond.InternalError(c)
		return
	}
	respond.OK(c, items)
}

// GetByID godoc
//
//	@Summary		Get campus by ID
//	@Tags			campuses
//	@Produce		json
//	@Param			id	path		string	true	"Campus ID"
//	@Success		200	{object}	Response
//	@Failure		400	{object}	respond.ErrorBody
//	@Failure		404	{object}	respond.ErrorBody
//	@Failure		500	{object}	respond.ErrorBody
//	@Router			/campuses/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	item, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if err.Error() == "invalid id" {
			respond.BadRequest(c, "invalid campus id")
			return
		}
		if errors.Is(err, ErrNotFound) {
			respond.NotFound(c, "campus not found")
			return
		}
		respond.InternalError(c)
		return
	}
	respond.OK(c, item)
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"module": "campus",
	})
}
