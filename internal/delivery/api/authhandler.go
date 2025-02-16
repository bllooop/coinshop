package api

import (
	"net/http"

	"github.com/bllooop/coinshop/internal/domain"
	logger "github.com/bllooop/coinshop/pkg/logging"
	"github.com/gin-gonic/gin"
)

func (h *Handler) SignUp(c *gin.Context) {
	logger.Log.Info().Msg("Получили запрос на создание пользователя")
	if c.Request.Method != http.MethodPost {
		newErrorResponse(c, http.StatusBadRequest, "Требуется запрос POST")
		logger.Log.Error().Msg("Требуется запрос POST")
		return
	}
	var input domain.User
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		logger.Log.Error().Err(err).Msg("")
		return
	}
	logger.Log.Debug().Msgf("Успешно прочитаны никнейм: %s, пароль: %s", input.UserName, input.Password)
	defaultCoins := 1000
	input.Coins = &defaultCoins
	id, err := h.Usecases.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		logger.Log.Error().Err(err).Msg("")
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
	logger.Log.Info().Msg("Создали пользователя")

}

func (h *Handler) SignIn(c *gin.Context) {
	logger.Log.Info().Msg("Получили запрос на авторизацию пользователя")

	var input domain.SignInInput
	var inputCreate domain.User
	if c.Request.Method != http.MethodPost {
		logger.Log.Error().Msg("Требуется запрос POST")
		newErrorResponse(c, http.StatusBadRequest, "Требуется запрос POST")
		return
	}
	if err := c.BindJSON(&input); err != nil {
		logger.Log.Error().Err(err).Msg("")
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	logger.Log.Debug().Msgf("Успешно прочитаны никнейм: %s, пароль: %s", input.UserName, input.Password)
	user, err := h.Usecases.Authorization.SignUser(input.UserName, input.Password)
	if err != nil {
		if err.Error() == "пользователь не найден" {
			logger.Log.Info().Msg("Создаем пользователя")
			inputCreate.UserName = input.UserName
			inputCreate.Password = input.Password
			defaultCoins := 1000
			inputCreate.Coins = &defaultCoins

			id, err2 := h.Usecases.Authorization.CreateUser(inputCreate)
			if err2 != nil {
				logger.Log.Error().Err(err2).Msg("")
				newErrorResponse(c, http.StatusInternalServerError, err2.Error())
				return
			}

			token, tokenErr := h.Usecases.Authorization.GenerateToken(id)
			if tokenErr != nil {
				logger.Log.Error().Err(tokenErr).Msg("")
				newErrorResponse(c, http.StatusInternalServerError, "Ошибка создания токена: "+tokenErr.Error())
				return
			}

			c.JSON(http.StatusOK, map[string]interface{}{
				"id":    id,
				"token": token,
			})
			logger.Log.Info().Msg("Получили токен")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, "Ошибка авторизации: "+err.Error())
		return
	}

	token, err := h.Usecases.Authorization.GenerateToken(user.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Ошибка создания токена: "+err.Error())
		logger.Log.Error().Err(err).Msg("")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
	logger.Log.Info().Msg("Получили токен")
}
