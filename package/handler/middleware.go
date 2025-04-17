package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	refreshTokenHeader  = "Refresh-Token"
	userCtxKey          = "username"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.Request.Header.Get(authorizationHeader)
	if header == "" {
		NewErrorResponse(c, http.StatusUnauthorized, "No authorization header")
		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 {
		NewErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header")
		return
	}

	username, _, err := h.services.Auth.ParseToken(headerParts[1])
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(userCtxKey, username)
}

func (h *Handler) userRefreshTokenCheck(c *gin.Context) {
	header := c.Request.Header.Get(authorizationHeader)
	if header == "" {
		NewErrorResponse(c, http.StatusUnauthorized, "No authorization header")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 1 {
		NewErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header")
		return
	}

	refreshHeader := c.Request.Header.Get(refreshTokenHeader)
	if refreshHeader == "" {
		NewErrorResponse(c, http.StatusUnauthorized, "No refresh token header")
		return
	}
	refreshParts := strings.Split(refreshHeader, " ")
	if len(refreshParts) != 1 {
		NewErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token header")
		return
	}
	maybeAccessToken, err := h.services.Auth.ParseRefreshToken(refreshParts[0])
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	if maybeAccessToken != headerParts[0] {
		NewErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token")
	}

	username, ip, err := h.services.Auth.ParseToken(headerParts[0])
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
	}
	if ip != c.ClientIP() {
		err := h.services.Auth.SendEmail(c.ClientIP(), username)
		if err != nil {
			NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		}
	}
	isTrueToken, err := h.services.Auth.VerifyToken(refreshParts[0], username)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	if isTrueToken {
		c.Set(userCtxKey, username)
	} else {
		NewErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token")
	}

}
