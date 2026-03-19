package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nexoai/backend/services"
	"github.com/stripe/stripe-go/v82"
)

type CheckoutHandler struct {
	stripe *services.StripeService
	email  *services.EmailService
	crm    *services.CRMService
}

func NewCheckoutHandler(s *services.StripeService, e *services.EmailService, c *services.CRMService) *CheckoutHandler {
	return &CheckoutHandler{stripe: s, email: e, crm: c}
}

func (h *CheckoutHandler) CreateCheckout(c *gin.Context) {
	var req services.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	if req.Email == "" || req.Plan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email y plan son requeridos"})
		return
	}

	if req.Plan != "monthly" && req.Plan != "annual" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan debe ser 'monthly' o 'annual'"})
		return
	}

	session, err := h.stripe.CreateCheckoutSession(req)
	if err != nil {
		log.Printf("[Checkout] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando sesión de pago"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *CheckoutHandler) HandleWebhook(c *gin.Context) {
	const maxBodyBytes = 65536
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodyBytes)

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("[Webhook] Error leyendo body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo leer el body"})
		return
	}

	sigHeader := c.GetHeader("Stripe-Signature")
	if sigHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Header Stripe-Signature requerido"})
		return
	}

	event, err := h.stripe.ConstructEvent(payload, sigHeader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Firma de webhook inválida"})
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		h.onCheckoutCompleted(event)
	default:
		log.Printf("[Webhook] Evento no manejado: %s", event.Type)
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

func (h *CheckoutHandler) onCheckoutCompleted(event stripe.Event) {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		log.Printf("[Webhook] Error parseando checkout.session: %v", err)
		return
	}

	customerEmail := session.CustomerEmail
	if customerEmail == "" && session.CustomerDetails != nil {
		customerEmail = session.CustomerDetails.Email
	}

	if customerEmail == "" {
		log.Println("[Webhook] checkout.session.completed sin email de cliente")
		return
	}

	log.Printf("[Webhook] Pago completado para %s (session: %s)", customerEmail, session.ID)

	go func() {
		if err := h.crm.CreateContact(services.ContactRequest{Email: customerEmail}); err != nil {
			log.Printf("[Webhook] Error creando contacto en HubSpot: %v", err)
		}

		if err := h.email.SendWelcomeEmail(services.WelcomeEmailRequest{
			To:   customerEmail,
			Name: customerEmail,
		}); err != nil {
			log.Printf("[Webhook] Error enviando email de bienvenida: %v", err)
		}
	}()
}
