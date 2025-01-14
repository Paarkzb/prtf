package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	userCtx = "userId"
)

type userIdentityInput struct {
	Auth   bool      `json:"auth" binding:"required"`
	UserID uuid.UUID `json:"userID" binding:"required"`
}

func (h *Handler) userIdentity(c *gin.Context) {
	var input userIdentityInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, input.UserID)
}

func getUserId(c *gin.Context) (uuid.UUID, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return uuid.Nil, errors.New("user id not found")
	}

	idUUID, ok := id.(uuid.UUID)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
		return uuid.Nil, errors.New("user id is of invalid type")
	}

	return idUUID, nil
}
