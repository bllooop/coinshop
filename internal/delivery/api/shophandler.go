package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/bllooop/coinshop/internal/domain"
	logger "github.com/bllooop/coinshop/pkg/logging"
	"github.com/gin-gonic/gin"
)

func (h *Handler) SendCoin(c *gin.Context) {
	logger.Log.Info().Msg("Получен запрос на отправку монет")
	if c.Request.Method != http.MethodPost {
		logger.Log.Error().Msg("Требуется запрос POST")
		newErrorResponse(c, http.StatusBadRequest, "Требуется запрос POST")
		return
	}
	userId, err := getUserId(c)
	if err != nil {
		logger.Log.Error().Err(err).Msg("")
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input domain.Transactions
	if err = c.BindJSON(&input); err != nil {
		logger.Log.Error().Err(err).Msg("")
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if input.DestinationUsername == "" || input.Amount < 0 {
		logger.Log.Error().Msg("Значения получателя и суммы не могут быть отрицательными или пустыми")
		newErrorResponse(c, http.StatusBadRequest, "Значения получателя и суммы не могут быть отрицательными или пустыми")
		return
	}
	logger.Log.Debug().Msgf("Успешно прочитано id  %v", userId)
	logger.Log.Debug().Msgf("Успешно прочитаны никнейм получателя: %s, количество: %v", input.DestinationUsername, input.Amount)
	id, err := h.Usecases.Shop.SendCoin(userId, input)
	if err != nil {
		logger.Log.Error().Err(err).Msg("")
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Log.Info().Msg("Получен ответ на отправку монет")

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) BuyItem(c *gin.Context) {
	logger.Log.Info().Msg("Получили запрос на покупку товара")
	if c.Request.Method != http.MethodPut {
		logger.Log.Error().Msg("Требуется запрос PUT")
		newErrorResponse(c, http.StatusBadRequest, "Требуется запрос PUT")
		return
	}
	userId, err := getUserId(c)
	if err != nil {
		logger.Log.Error().Err(err).Msg("")
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	name := c.Param("item")
	logger.Log.Debug().Msgf("Успешно прочитаны название предмета %s и id  %v", name, userId)

	id, err := h.Usecases.Shop.BuyItem(userId, name)
	if err != nil {
		logger.Log.Error().Err(err).Msg("")
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Log.Info().Msg("Получен ответ на покупку товара")

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetInfo(c *gin.Context) {
	logger.Log.Info().Msg("Получили запрос на информацию о пользователе")
	if c.Request.Method != http.MethodGet {
		logger.Log.Error().Msg("Требуется запрос GET")
		newErrorResponse(c, http.StatusBadRequest, "Требуется запрос GET")
		return
	}
	userId, err := getUserId(c)
	if err != nil {
		logger.Log.Error().Err(err).Msg("")
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Log.Debug().Msgf("Успешно прочитано id  %v", userId)
	lists, err := h.Usecases.Shop.GetUserSummary(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Error().Err(err).Msg("Пользователь не найден")
			newErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
			return
		}
		logger.Log.Error().Err(err).Msg("")
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Log.Info().Msg("Получен ответ на запрос информации о пользователе")

	c.JSON(http.StatusOK, lists)
}
