package handler

import (
	"prtf/pkg/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	mux := gin.Default()

	mux.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8080"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept", "Accept-Encoding", "User-Agent", "Cache-Control", "Pragma", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	auth := mux.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := mux.Group("/api", h.userIdentity)
	{
		quiz := api.Group("/quiz")
		{
			quiz.POST("", h.saveQuiz)
			quiz.GET("", h.getAllQuiz)
			quiz.GET("/:id", h.getQuizById)
			quiz.PUT("/:id", h.updateQuiz)
			quiz.PATCH("/:id", h.updateQuiz)
			quiz.DELETE("/:id", h.deleteQuiz)

		}
	}

	return mux
}
