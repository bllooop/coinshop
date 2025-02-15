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
		newErrorResponse(c, http.StatusBadRequest, "Требуется запрос POST")
		return
	}
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Log.Debug().Msgf("Успешно прочитано id  %v", userId)
	var input domain.Transactions
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if input.DestinationUsername == "" || input.Amount < 0 {
		newErrorResponse(c, http.StatusBadRequest, "Значения получателя и суммы не могут быть отрицательными или пустыми")
		return
	}
	id, err := h.Usecases.Shop.SendCoin(userId, input)
	if err != nil {
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
		newErrorResponse(c, http.StatusBadRequest, "Требуется запрос PUT")
		return
	}
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	name := c.Param("item")
	logger.Log.Debug().Msgf("Успешно прочитаны название предмета %s и id  %v", name, userId)

	id, err := h.Usecases.Shop.BuyItem(userId, name)
	if err != nil {
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
		newErrorResponse(c, http.StatusBadRequest, "Требуется запрос GET")
		return
	}
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Log.Debug().Msgf("Успешно прочитано id  %v", userId)
	lists, err := h.Usecases.Shop.GetUserSummary(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newErrorResponse(c, http.StatusNotFound, "User not found")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Log.Info().Msg("Получен ответ на запрос информации о пользователе")

	c.JSON(http.StatusOK, lists)
}
