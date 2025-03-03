package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) getChannelRecordings(c *gin.Context) {
	id := c.Param("id")
	channelID, err := uuid.Parse(id)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid id params")
		return
	}

	recordings, err := h.streamService.GetChannelRecordings(c, channelID)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, "failed to get recordings")
		return
	}

	c.JSON(http.StatusOK, recordings)
}

func (h *Handler) getRecordingById(c *gin.Context) {
	id := c.Param("id")
	recordingId, err := uuid.Parse(id)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid id params")
		return
	}

	recording, err := h.streamService.GetRecordingById(c, recordingId)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, "failed to get recording")
		return
	}

	c.JSON(http.StatusOK, recording)
}

func (h *Handler) authStream(c *gin.Context) {
	streamKey := c.Query("name")

	h.log.Infow("auth stream", "streamKey", streamKey)

	if streamKey == "" {
		h.newErrorResponse(c, http.StatusBadRequest, "stream key required")
		return
	}

	channel, err := h.streamService.ValidateStreamToken(c, streamKey)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, "failed to auth stream")
		return
	}

	h.log.Infow("streamKey valid", "channelID", channel.ID)

	_, err = h.streamService.StartStream(c, channel.ID)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, "failed to start stream")
		return
	}

	//

	c.Writer.Header().Add("Location", channel.ChannelName)

	c.Status(http.StatusMovedPermanently)
}

func (h *Handler) endStream(c *gin.Context) {
	streamKey := c.Query("name")

	h.log.Infow("end stream", "streamKey", streamKey)

	if streamKey == "" {
		h.newErrorResponse(c, http.StatusBadRequest, "stream key required")
		return
	}

	channel, err := h.streamService.ValidateStreamToken(c, streamKey)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, "failed to auth stream")
		return
	}

	h.log.Infow("streamKey valid", "channelID", channel.ID)

	_, err = h.streamService.EndStream(c, channel)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, "failed to end stream")
		return
	}

	c.Status(http.StatusOK)
}
