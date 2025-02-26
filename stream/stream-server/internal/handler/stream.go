package handler

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) listChannelRecordings(c *gin.Context) {
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

	cmd := exec.Command("sh", "/var/scripts/save_record.sh", channel.ChannelName)
	recordPath, err := cmd.Output()
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, "failed to save stream record")
		return
	}

	_, err = h.streamService.EndStream(c, channel.ID, string(recordPath))
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, "failed to end stream")
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) listStreams(c *gin.Context) {

	// c.JSON(http.StatusOK, activeStreams)
	c.Status(http.StatusOK)
}
