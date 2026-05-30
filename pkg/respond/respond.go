package respond

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Error string `json:"error"`
}

func JSON(c *gin.Context, status int, body any) {
	c.JSON(status, body)
}

func OK(c *gin.Context, body any) {
	c.JSON(http.StatusOK, body)
}

func Created(c *gin.Context, body any) {
	c.JSON(http.StatusCreated, body)
}

func Error(c *gin.Context, status int, msg string) {
	c.JSON(status, ErrorBody{Error: msg})
}

func BadRequest(c *gin.Context, msg string) {
	Error(c, http.StatusBadRequest, msg)
}

func Unauthorized(c *gin.Context, msg string) {
	Error(c, http.StatusUnauthorized, msg)
}

func NotFound(c *gin.Context, msg string) {
	Error(c, http.StatusNotFound, msg)
}

func Conflict(c *gin.Context, msg string) {
	Error(c, http.StatusConflict, msg)
}

func InternalError(c *gin.Context) {
	Error(c, http.StatusInternalServerError, "internal server error")
}
