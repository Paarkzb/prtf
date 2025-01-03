package handler

import (
	"errors"
	"net/http"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/repository"
	"sso/internal/services/authservice"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

		h.newErrorResponse(c, http.StatusBadRequest, errors.New("failed to sign-up"))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"userID": id,
	})
}

func (h *Handler) signIn(c *gin.Context) {
	var input models.UserInput

	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	tokens, err := h.authService.SignIn(c, input.Username, input.Password)
	if err != nil {
		if errors.Is(err, authservice.ErrInvalidCredentials) {
			h.newErrorResponse(c, http.StatusInternalServerError, errors.New("invalid username or password"))
			return
		}

		h.newErrorResponse(c, http.StatusInternalServerError, errors.New("failed to login"))
		return
	}

	tokenClaims, err := jwt.ParseToken(tokens.AccessToken)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	user, err := h.authService.GetUserByUserID(c, uuid.MustParse(tokenClaims["uid"].(string)))
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

type isAdminInput struct {
	UserID uuid.UUID `json:"userID" binding:"required"`
}

func (h *Handler) isAdmin(c *gin.Context) {
	var input isAdminInput

	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	isAdmin, err := h.authService.IsAdmin(c, input.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			h.newErrorResponse(c, http.StatusInternalServerError, errors.New("user not found"))
			return
		}
		h.newErrorResponse(c, http.StatusInternalServerError, errors.New("failed to check admin status"))
		return

	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"isAdmin": isAdmin,
	})
}

type userIdentityInput struct {
	AccessToken string `json:"accessToken"`
}

func (h *Handler) userIdentity(c *gin.Context) {
	var input userIdentityInput

	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	auth, userID, err := h.authService.UserIdentity(c, input.AccessToken)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, errors.New("authentication failed"))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"auth":   auth,
		"userID": userID,
	})
}

type refreshInput struct {
	UserID       uuid.UUID `json:"userID"`
	RefreshToken string    `json:"refreshToken"`
}

func (h *Handler) refresh(c *gin.Context) {
	var input refreshInput

	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	tokens, err := h.authService.Refresh(c, input.UserID, input.RefreshToken)
	if err != nil {
		if errors.Is(err, authservice.ErrInvalidCredentials) {
			h.newErrorResponse(c, http.StatusInternalServerError, errors.New("invalid userID or refresh token"))
			return
		}
		h.newErrorResponse(c, http.StatusInternalServerError, errors.New("failed to refresh token"))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
	})
}
