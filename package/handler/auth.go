package handler

import (
	"REST_JWT"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var input REST_JWT.User
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.CreateUser(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}

type signInRequest struct {
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input REST_JWT.User
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	accesstoken, err := h.services.GenerateToken(input.Username, input.Password, c.ClientIP())
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	refreshToken, err := h.services.GenerateRefreshToken(accesstoken)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accesstoken":  accesstoken,
		"refreshtoken": refreshToken,
	})
	err = h.services.SaveRefreshToken(refreshToken, input.Username)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
}

func (h *Handler) signOut(c *gin.Context) {
	return

}

func (h *Handler) refreshToken(c *gin.Context) {
	h.userRefreshTokenCheck(c)
	nameRaw, exists := c.Get(userCtxKey)
	if !exists {
		NewErrorResponse(c, http.StatusUnauthorized, "user not authorized")
		return
	}
	name, ok := nameRaw.(string)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "invalid user context")
		return
	}
	accesstoken, err := h.services.GenerateNewToken(name, c.ClientIP())
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	refreshToken, err := h.services.GenerateRefreshToken(accesstoken)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"accesstoken":  accesstoken,
		"refreshtoken": refreshToken,
	})
	err = h.services.SaveRefreshToken(refreshToken, name)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
}
