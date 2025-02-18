package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"videostream/internal/domain/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) saveChannel(c *gin.Context) {
	userID, err := h.getUserId(c)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var channel models.Channel
	channel.UserID = userID

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
		h.newErrorResponse(c, http.StatusInternalServerError, "failed to save channel")
		return
	}

	c.JSON(http.StatusOK, channels)
}

type Recording struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Date     string `json:"date"`
	Duration string `json:"duration"`
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

func isValidStreamKey(key string) bool {
	return key == "test123"
}

func (h *Handler) authStreamHandler(c *gin.Context) {
	key := c.Query("name")

	h.log.Infow(key, " started stream")

	if key == "" {
		h.newErrorResponse(c, http.StatusBadRequest, "stream key required")
		return
	}

	if isValidStreamKey(key) {
		// mu.Lock()
		// activeStreams[key] = true
		// mu.Unlock()

		c.Status(http.StatusOK)
		return
	}

	h.newErrorResponse(c, http.StatusUnauthorized, "invalid stream key")
}

func (h *Handler) listStreamsHandler(c *gin.Context) {

	// cached, err := rdb.Get(context.Background(), "streams").Bytes()
	// if err == nil {
	// 	w.Write(cached)
	// 	return
	// }

	// mu.Lock()
	// rdb.Set(context.Background(), "streams", activeStreams, 10*time.Second)

	// defer mu.Unlock()

	// c.JSON(http.StatusOK, activeStreams)
	c.Status(http.StatusOK)
}
