package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	StripeSecretKey  string
	StripeWebhookSec string
	MonthlyPriceID   string
	AnnualPriceID    string
	SuccessURL       string
	CancelURL        string
	ResendAPIKey     string
	HubSpotAPIKey    string
	TruoraAPIKey     string
	TruoraFlowID     string
	TruoraRedirectURL string
	AllowedOrigins   string
	FrameAncestors   string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[Config] No se encontró archivo .env, usando variables de entorno del sistema")
	}

	return &Config{
		Port:             getEnv("PORT", "8080"),
		StripeSecretKey:  getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookSec: getEnv("STRIPE_WEBHOOK_SECRET", ""),
		MonthlyPriceID:   getEnv("STRIPE_MONTHLY_PRICE_ID", ""),
		AnnualPriceID:    getEnv("STRIPE_ANNUAL_PRICE_ID", ""),
		SuccessURL:       getEnv("SUCCESS_URL", "http://localhost:5173/success"),
		CancelURL:        getEnv("CANCEL_URL", "http://localhost:5173/cancel"),
		ResendAPIKey:     getEnv("RESEND_API_KEY", ""),
		HubSpotAPIKey:    getEnv("HUBSPOT_API_KEY", ""),
		TruoraAPIKey:      getEnv("TRUORA_API_KEY", ""),
		TruoraFlowID:      getEnv("TRUORA_FLOW_ID", ""),
		TruoraRedirectURL: getEnv("TRUORA_REDIRECT_URL", "http://localhost:5173/verify"),
		AllowedOrigins:   getEnv("ALLOWED_ORIGINS", "*"),
		FrameAncestors:   getEnv("FRAME_ANCESTORS", "*"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
