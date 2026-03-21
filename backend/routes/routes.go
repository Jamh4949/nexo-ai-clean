package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nexoai/backend/handlers"
)

func Setup(r *gin.Engine, checkout *handlers.CheckoutHandler, truora *handlers.TruoraHandler) {
	api := r.Group("/api/v1")
	{
		api.GET("/health", handlers.HealthCheck)
		api.POST("/checkout", checkout.CreateCheckout)
		api.POST("/webhook/stripe", checkout.HandleWebhook)
	}

	r.POST("/api/truora/generate-token", truora.GenerateToken)
}
