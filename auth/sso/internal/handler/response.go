package handler

import (
	"sso/internal/lib/logger/sl"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func (h *Handler) newErrorResponse(c *gin.Context, statusCode int, err error) {
	h.log.Error("failed to generate password hash", sl.Err(err))
	c.AbortWithStatusJSON(statusCode, errorResponse{Message: err.Error()})
}
