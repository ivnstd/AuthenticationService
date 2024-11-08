package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authenticationHeader = "Authorization"
	userCtx              = "userID"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authenticationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "Empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid auth header")
		return
	}

	userID, err := h.services.Auth.ParseAccessToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userID)
}
