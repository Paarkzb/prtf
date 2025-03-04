package handler

import (
	"errors"
	"net/http"
	"strings"
	"videostream/internal/lib/jwt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	userCtx    = "userId"
	channelCtx = "channelId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := strings.Split(c.Request.Header["Authorization"][0], " ")

	h.log.Infow("userIdentity", "header", header)

	if len(header) < 2 {
		h.newErrorResponse(c, http.StatusBadRequest, "header is empty")
		return
	}

	token := header[1]
	if token == "" {
		h.newErrorResponse(c, http.StatusBadRequest, "token is empty")
		return
	}

	claims, err := jwt.ParseToken(token)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.log.Infow("userIdentity", "claims", claims)

	c.Set(userCtx, claims["uid"])
}

func (h *Handler) channelIdentity(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	channel, err := h.streamService.GetChannelByUserId(c, userId)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "failed to get channel data")
		return
	}

	c.Set(channelCtx, channel.ID)
}

func (h *Handler) getUserId(c *gin.Context) (uuid.UUID, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		h.newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return uuid.Nil, errors.New("user id not found")
	}

	h.log.Infow("getUserId", "id", id)

	idUUID, err := uuid.Parse(id.(string))
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return uuid.Nil, errors.New("user id is of invalid type")
	}

	return idUUID, nil
}
