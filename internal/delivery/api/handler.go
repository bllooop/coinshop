package api

import (
	"github.com/bllooop/coinshop/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Usecases *usecase.Usecase
}

func NewHandler(usecases *usecase.Usecase) *Handler {
	return &Handler{Usecases: usecases}
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
			auth.POST("sign-up", h.SignUp)
			auth.POST("sign-in", h.SignIn)
		}
		authorized := api.Group("/", h.authIdentity)
		//authorized.Use(h.AuthMiddleware)
		{
			authorized.POST("/sendCoin", h.SendCoin)
			authorized.GET("/info", h.GetInfo)
			authorized.PUT("/buy/:item", h.BuyItem)
		}
	}
	return router
}
