package handler

import (
	"net/http"
	"videostream/internal/services/streamservice"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Handler struct {
	log           *zap.SugaredLogger
	streamService *streamservice.StreamService
}

func NewHandler(log *zap.SugaredLogger, streamService *streamservice.StreamService) *Handler {
	return &Handler{
		log:           log,
		streamService: streamService,
	}
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	mux := gin.Default()

	mux.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:80", "http://localhost:5173", "http://localhost:8090"},
		// AllowAllOrigins:  true,
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept", "Accept-Encoding", "User-Agent", "Cache-Control", "Pragma", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.GET("/metrics", prometheusHandler())

	mux.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	mux.GET("/stream/start", h.authStream)
	mux.GET("/stream/end", h.endStream)

	streamApi := mux.Group("/streams", h.userIdentity)
	{
		streamApi.GET("", h.listStreams)
	}

	// mux.Handle("/vod/", http.StripPrefix("/vod/", http.FileServer(http.Dir("/var/vod"))))

	channelApi := mux.Group("/channels", h.userIdentity)
	{
		channelApi.POST("", h.saveChannel)
		channelApi.GET("", h.getAllChannels)
		channelApi.GET("/:id", h.getChannelById)
		channelApi.GET("/:id/recordings", h.listChannelRecordings)

		channelApi.GET("/user", h.getChannelByUserId)
	}

	return mux
}
