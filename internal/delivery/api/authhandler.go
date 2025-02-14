package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/bllooop/coinshop/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input domain.User
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	defaultCoins := 1000
	input.Coins = &defaultCoins
	id, err := h.usecases.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) signIn(c *gin.Context) {
	var input domain.SignInInput
	var inputCreate domain.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	user, err := h.usecases.Authorization.SignUser(input.UserName, input.Password)
	if err == nil {
		token, err := h.usecases.Authorization.GenerateToken(user.Id)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, "Failed to generate token: "+err.Error())
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"token": token,
		})
		return
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		defaultCoins := 1000
		inputCreate.Coins = &defaultCoins
		id, err := h.usecases.Authorization.CreateUser(inputCreate)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		token, tokenErr := h.usecases.Authorization.GenerateToken(id)
		if tokenErr != nil {
			newErrorResponse(c, http.StatusInternalServerError, "Failed to generate token: "+tokenErr.Error())
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"id":    id,
			"token": token,
		})
		return
	}
	newErrorResponse(c, http.StatusInternalServerError, "Error checking user: "+err.Error())
}
