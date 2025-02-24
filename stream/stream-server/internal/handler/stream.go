package handler

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type Recording struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Date     string `json:"date"`
	Duration string `json:"duration"`
}

func (h *Handler) startStream(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	channel, err := h.streamService.GetChannelByUserId(c, userId)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "channel not found")
		return
	}

	// streamToken, err := h.streamService.GenerateStreamToken(c, channel)
	// if err != nil {
	// 	h.newErrorResponse(c, http.StatusInternalServerError, "failed to start stream")
	// 	return
	// }

	c.JSON(http.StatusOK, channel.ChannelToken)
}

func (h *Handler) listRecordingsHandler(c *gin.Context) {
	files, err := os.ReadDir("/var/vod")
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, "failed to read recordings")
		return
	}

	var recordings []Recording
	for _, f := range files {
		fileInfo, err := f.Info()
		if err != nil {
			h.newErrorResponse(c, http.StatusInternalServerError, "failed to read recordings info")
			return
		}

		if filepath.Ext(f.Name()) == ".mp4" {
			recordings = append(recordings, Recording{
				Name: f.Name(),
				Path: f.Name(),
				Date: fileInfo.ModTime().Format("2006-01-02 15:14"),
			})
		}
	}

	c.JSON(http.StatusOK, recordings)
}

func (h *Handler) authStreamHandler(c *gin.Context) {
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

func (h *Handler) listStreamsHandler(c *gin.Context) {

	// c.JSON(http.StatusOK, activeStreams)
	c.Status(http.StatusOK)
}
