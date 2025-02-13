package api

import (
	"net/http"

	"github.com/bllooop/coinshop/internal/domain"
	logger "github.com/bllooop/coinshop/pkg/logging"
	"github.com/gin-gonic/gin"
)

func (h *Handler) sendCoin(c *gin.Context) {
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
	//userId := 1
	var input domain.Transactions
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if input.Destination < 1 || input.Amount < 0 {
		newErrorResponse(c, http.StatusBadRequest, "Значения получателя и суммы не могут быть отрицательными")
		return
	}
	id, err := h.usecases.Shop.SendCoin(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Log.Info().Msg("Получен ответ на отправку монет")

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) buyItem(c *gin.Context) {
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
	name := c.Param("name")
	logger.Log.Debug().Msgf("Успешно прочитаны name: %s", name)

	id, err := h.usecases.Shop.BuyItem(userId, name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	logger.Log.Info().Msg("Получен ответ на покупку товара")

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
