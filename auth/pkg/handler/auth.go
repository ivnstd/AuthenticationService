package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ivnstd/AuthenticationService/auth/models"
)

func (h *Handler) signUp(c *gin.Context) {
	var input models.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid sign-up format")
		return
	}

	if err := h.services.Auth.CreateUser(input); err != nil {
		switch {
		case strings.Contains(err.Error(), "uni_users_username"):
			newErrorResponse(c, http.StatusConflict, "User with this username already exists")
		case strings.Contains(err.Error(), "uni_users_email"):
			newErrorResponse(c, http.StatusConflict, "User with this email already exists")
		default:
			newErrorResponse(c, http.StatusInternalServerError, "Failed to sign-up")
		}
		return
	}

	newSuccessResponse(c, http.StatusOK, "message", "Successful sign-up")
}

func (h *Handler) signIn(c *gin.Context) {
	var input models.SignInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid sign-in format")
		return
	}

	user, err := h.services.Auth.GetUserByUsername(input.Username, input.Password)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			newErrorResponse(c, http.StatusInternalServerError, "Invalid username or password", err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, "Failed to sign-in", err.Error())
		}
		return
	}

	accessToken, err := h.services.Auth.GenerateAccessToken(user.ID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	refreshToken, err := h.services.Auth.GenerateRefreshToken(user.ID, c.ClientIP())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	newSuccessResponse(c, http.StatusOK, "tokens", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) info(c *gin.Context) {
	id, _ := c.Get(userCtx)
	newSuccessResponse(c, http.StatusOK, "ID", id)
}

func (h *Handler) refresh(c *gin.Context) {
	var input models.RefreshInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	newAccessToken, newRefreshToken, err := h.services.Auth.RefreshTokens(input.AccessToken, input.RefreshToken, c.ClientIP())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to refresh tokens")
		return
	}
	newSuccessResponse(c, http.StatusOK, "tokens", gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func (h *Handler) logout(c *gin.Context) {
	var input models.LogoutInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	if err := h.services.Auth.RevokeRefreshToken(input.RefreshToken); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to log out")
		return
	}

	newSuccessResponse(c, http.StatusOK, "message", "Successfully logged out")
}
