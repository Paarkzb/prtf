package handler

import (
	"errors"
	"net/http"
	"prtf/internal/lib/jwt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	userCtx = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := strings.Split(c.Request.Header["Authorization"][0], " ")

	if len(header) < 2 {
		newErrorResponse(c, http.StatusBadRequest, "header is empty")
		return
	}

	token := header[1]

	if token == "" {
		newErrorResponse(c, http.StatusBadRequest, "token is empty")
		return
	}

	claims, err := jwt.ParseToken(token)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("userIdentity. claims: %v", claims)

	c.Set(userCtx, claims["uid"])
}

func getUserId(c *gin.Context) (uuid.UUID, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return uuid.Nil, errors.New("user id not found")
	}

	logrus.Infof("getUserId. id: %s", id)

	idUUID, err := uuid.Parse(id.(string))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return uuid.Nil, errors.New("user id is of invalid type")
	}

	logrus.Infof("valid auth. userID: %s", idUUID)

	return idUUID, nil
}
