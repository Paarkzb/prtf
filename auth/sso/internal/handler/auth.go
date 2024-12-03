package handler

import (
	"errors"
	"net/http"
	"sso/internal/domain/models"
	"sso/internal/services/authservice"

	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input models.UserInput

	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.authService.SignUp(c, input.Username, input.Email, input.Password)
	if err != nil {
		if errors.Is(err, authservice.ErrInvalidCredentials) {
			h.newErrorResponse(c, http.StatusBadRequest, errors.New("invalid username or password"))
			return
		}

		h.newErrorResponse(c, http.StatusBadRequest, errors.New("failed to login"))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"userID": id,
	})
}

func (h *Handler) signIn(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
func (h *Handler) isAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
func (h *Handler) userIdentity(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
func (h *Handler) refresh(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
