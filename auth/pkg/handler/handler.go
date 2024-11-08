package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ivnstd/AuthenticationService/auth/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.refresh)
		auth.POST("/logout", h.logout)
	}

	api := router.Group("/api", h.userIdentity)
	{
		api.GET("info", h.info)
	}

	return router
}
