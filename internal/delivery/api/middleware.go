package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) authIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "Пустой заголовок авторизации")
		c.Abort()
		return
	}
	headerSplit := strings.Split(header, " ")
	if len(headerSplit) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "Некорректный ввод токена")
		c.Abort()
		return
	}
	userId, err := h.usecases.Authorization.ParseToken(headerSplit[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}
	c.Set(userCtx, userId)
}
func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("ID пользователя не найдено")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("ID пользователя некорректного типа данных")
	}

	return idInt, nil
}
