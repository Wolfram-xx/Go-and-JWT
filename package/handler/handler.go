package handler

import (
	"REST_JWT/package/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh-token", h.refreshToken)
		auth.POST("/logout", h.signOut)
	}
	api := router.Group("/api", h.userIdentity)
	{
		api.GET("/someImportantInformation", h.getInformation)
	}
	return router
}
