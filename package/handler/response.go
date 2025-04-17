package handler

import (
	"github.com/gin-gonic/gin"
)

type Error struct {
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, Error{message})
}
