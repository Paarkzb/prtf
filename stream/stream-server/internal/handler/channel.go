package handler

import (
	"net/http"
	"videostream/internal/domain/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) saveChannel(c *gin.Context) {
	userID, err := h.getUserId(c)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var channel models.Channel
	channel.RfUserID = userID

	_, err = h.streamService.SaveChannel(c, channel)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, "failed to save channel")
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) getAllChannels(c *gin.Context) {

	channels, err := h.streamService.GetAllChannels(c)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, "failed to get channels")
		return
	}

	c.JSON(http.StatusOK, channels)
}

func (h *Handler) getChannelById(c *gin.Context) {
	// userId, err := h.getUserId(c)
	// if err != nil {
	// h.newErrorResponse(c, http.StatusBadRequest, err.Error())
	// 	return
	// }

	id := c.Param("id")
	channelID, err := uuid.Parse(id)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid id params")
		return
	}

	channelData, err := h.streamService.GetChannelById(c, channelID)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, "failed to get channel data")
		return
	}

	c.JSON(http.StatusOK, channelData)
}

func (h *Handler) getChannelByUserId(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	channelData, err := h.streamService.GetChannelByUserId(c, userId)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, "failed to get channel data")
		return
	}

	c.JSON(http.StatusOK, channelData)
}
