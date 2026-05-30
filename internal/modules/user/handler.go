package user

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/truongle2004/campus_marketplace/pkg/auth"
	"github.com/truongle2004/campus_marketplace/pkg/respond"
)

type Handler struct {
	service UserService
}

func NewHandler(service UserService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) authUser(c *gin.Context) (*auth.User, bool) {
	user, ok := auth.UserFromContext(c.Request.Context())
	if !ok {
		respond.Unauthorized(c, "unauthorized")
		return nil, false
	}
	return user, true
}

// Bootstrap godoc
//
//	@Summary		Create app profile after Clerk sign-up
//	@Description	Links Clerk identity to a campus profile (first-time setup).
//	@Tags			users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		BootstrapRequest	true	"Bootstrap payload"
//	@Success		201		{object}	ProfileResponse
//	@Failure		400		{object}	respond.ErrorBody
//	@Failure		401		{object}	respond.ErrorBody
//	@Failure		409		{object}	respond.ErrorBody
//	@Failure		500		{object}	respond.ErrorBody
//	@Router			/users/me/bootstrap [post]
func (h *Handler) Bootstrap(c *gin.Context) {
	authUser, ok := h.authUser(c)
	if !ok {
		return
	}

	var req BootstrapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.BadRequest(c, err.Error())
		return
	}

	profile, err := h.service.Bootstrap(c.Request.Context(), authUser, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrAlreadyExists):
			respond.Conflict(c, "profile already exists")
		case errors.Is(err, ErrInvalidCampus):
			respond.BadRequest(c, "invalid campus")
		case errors.Is(err, ErrEmailRequired):
			respond.BadRequest(c, "email missing from token")
		default:
			respond.InternalError(c)
		}
		return
	}

	respond.Created(c, profile)
}

// GetMe godoc
//
//	@Summary		View own profile
//	@Tags			users
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	ProfileResponse
//	@Failure		401	{object}	respond.ErrorBody
//	@Failure		404	{object}	respond.ErrorBody
//	@Failure		500	{object}	respond.ErrorBody
//	@Router			/users/me [get]
func (h *Handler) GetMe(c *gin.Context) {
	authUser, ok := h.authUser(c)
	if !ok {
		return
	}

	profile, err := h.service.GetProfile(c.Request.Context(), authUser.ClerkUserID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			respond.NotFound(c, "profile not found; bootstrap required")
			return
		}
		respond.InternalError(c)
		return
	}

	respond.OK(c, profile)
}

// UpdateMe godoc
//
//	@Summary		Edit own profile
//	@Tags			users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		UpdateProfileRequest	true	"Profile fields"
//	@Success		200		{object}	ProfileResponse
//	@Failure		400		{object}	respond.ErrorBody
//	@Failure		401		{object}	respond.ErrorBody
//	@Failure		404		{object}	respond.ErrorBody
//	@Failure		500		{object}	respond.ErrorBody
//	@Router			/users/me [patch]
func (h *Handler) UpdateMe(c *gin.Context) {
	authUser, ok := h.authUser(c)
	if !ok {
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.BadRequest(c, err.Error())
		return
	}
	if req.FullName == nil && req.AvatarURL == nil {
		respond.BadRequest(c, "at least one field is required")
		return
	}

	profile, err := h.service.UpdateProfile(c.Request.Context(), authUser.ClerkUserID, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			respond.NotFound(c, "profile not found")
			return
		}
		respond.InternalError(c)
		return
	}

	respond.OK(c, profile)
}

// DeactivateMe godoc
//
//	@Summary		Deactivate own account
//	@Tags			users
//	@Security		BearerAuth
//	@Success		204
//	@Failure		401	{object}	respond.ErrorBody
//	@Failure		404	{object}	respond.ErrorBody
//	@Failure		500	{object}	respond.ErrorBody
//	@Router			/users/me [delete]
func (h *Handler) DeactivateMe(c *gin.Context) {
	authUser, ok := h.authUser(c)
	if !ok {
		return
	}

	if err := h.service.Deactivate(c.Request.Context(), authUser.ClerkUserID); err != nil {
		if errors.Is(err, ErrNotFound) {
			respond.NotFound(c, "profile not found")
			return
		}
		respond.InternalError(c)
		return
	}

	c.Status(204)
}
