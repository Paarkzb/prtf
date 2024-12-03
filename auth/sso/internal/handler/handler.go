package handler

import (
	"log/slog"
	"sso/internal/services/authservice"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	log         *slog.Logger
	authService *authservice.AuthService
}

func NewHandler(log *slog.Logger, authService *authservice.AuthService) *Handler {
	return &Handler{
		log:         log,
		authService: authService,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	mux := gin.Default()

	mux.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept", "Accept-Encoding", "User-Agent", "Cache-Control", "Pragma", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.GET("/", h.getAllSongs)
	mux.POST("/", h.saveSong)

	return mux
}
