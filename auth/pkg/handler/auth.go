package handler

import (
	"net/http"

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
		newErrorResponse(c, http.StatusInternalServerError, "Failed to sign-up")
		return
	}

	newSuccessResponse(c, http.StatusOK, "message", "Successful sign-up")
}

func (h *Handler) signIn(c *gin.Context) {
	var input models.SignInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid sign-up format")
		return
	}

	token, err := h.services.Auth.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to sign-in")
		return
	}

	newSuccessResponse(c, http.StatusOK, "token", token)
}

func (h *Handler) info(c *gin.Context) {
	id, _ := c.Get(userCtx)
	newSuccessResponse(c, http.StatusOK, "ID", id)
}
