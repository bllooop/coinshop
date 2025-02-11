package api

import (
	"github.com/bllooop/coinshop/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecases *usecase.Usecase
}

func NewHandler(usecases *usecase.Usecase) *Handler {
	return &Handler{usecases: usecases}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("sign-up", h.signUp)
			auth.POST("sign-in", h.signIn)
		}
		authorized := api.Group("/", h.AuthMiddleware)
		//authorized.Use(h.AuthMiddleware)
		{
			authorized.POST("/sendCoin", h.sendCoin)
			//authorized.GET("/info", h.getInfo)
			authorized.PUT("/buy/:item", h.buyItem)
		}
	}
	return router
}
