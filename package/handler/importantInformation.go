package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getInformation(c *gin.Context) {
	name, _ := c.Get(userCtxKey)
	c.JSON(http.StatusOK, map[string]interface{}{
		"YourUsername": name,
	})
}
