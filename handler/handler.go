package handler

import (
	"github.com/gin-gonic/gin"
	"testTask/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{service: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/token")
	{
		auth.POST("/auth", h.auth)
		auth.POST("/refresh", h.refreshTokens)
	}

	router.POST("/sign-up", h.signUp)

	return router
}
