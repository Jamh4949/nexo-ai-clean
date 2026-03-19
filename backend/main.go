package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nexoai/backend/config"
	"github.com/nexoai/backend/handlers"
	"github.com/nexoai/backend/routes"
	"github.com/nexoai/backend/services"
)

func main() {
	cfg := config.Load()

	r := gin.Default()

	r.Use(corsMiddleware(cfg))
	r.Use(securityHeaders(cfg))

	stripeService := services.NewStripeService(cfg)
	emailService := services.NewEmailService(cfg)
	crmService := services.NewCRMService(cfg)

	checkoutHandler := handlers.NewCheckoutHandler(stripeService, emailService, crmService)

	routes.Setup(r, checkoutHandler)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("NexoAI API corriendo en %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}

func corsMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowed := cfg.AllowedOrigins

		if allowed == "*" {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			for _, o := range strings.Split(allowed, ",") {
				if strings.TrimSpace(o) == origin {
					c.Header("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func securityHeaders(cfg *config.Config) gin.HandlerFunc {
	csp := fmt.Sprintf("frame-ancestors %s", cfg.FrameAncestors)

	return func(c *gin.Context) {
		c.Header("Content-Security-Policy", csp)
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "ALLOWALL")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}
