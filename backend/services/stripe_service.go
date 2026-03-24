package services

import (
	"fmt"
	"log"

	"github.com/nexoai/backend/config"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/webhook"
)

type StripeService struct {
	secretKey      string
	webhookSec     string
	monthlyPriceID string
	annualPriceID  string
	successURL     string
	cancelURL      string
}

func NewStripeService(cfg *config.Config) *StripeService {
	stripe.Key = cfg.StripeSecretKey

	return &StripeService{
		secretKey:      cfg.StripeSecretKey,
		webhookSec:     cfg.StripeWebhookSec,
		monthlyPriceID: cfg.MonthlyPriceID,
		annualPriceID:  cfg.AnnualPriceID,
		successURL:     cfg.SuccessURL,
		cancelURL:      cfg.CancelURL,
	}
}

type CheckoutRequest struct {
	Plan  string `json:"plan"`  // "monthly" | "annual"
	Email string `json:"email"`
}

type CheckoutResponse struct {
	SessionID   string `json:"session_id"`
	CheckoutURL string `json:"checkout_url"`
}

func (s *StripeService) CreateCheckoutSession(req CheckoutRequest) (*CheckoutResponse, error) {
	if s.secretKey == "" {
		return nil, fmt.Errorf("STRIPE_SECRET_KEY no está configurado")
	}

	lineItem := s.buildLineItem(req.Plan)

	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(s.successURL + "?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(s.cancelURL),
		LineItems:  []*stripe.CheckoutSessionLineItemParams{lineItem},
	}

	if req.Email != "" {
		params.CustomerEmail = stripe.String(req.Email)
	}

	sess, err := session.New(params)
	if err != nil {
		log.Printf("[Stripe] Error creando checkout session: %v", err)
		return nil, fmt.Errorf("error creando sesión de Stripe: %w", err)
	}

	log.Printf("[Stripe] Checkout session creada: %s para %s (plan: %s)", sess.ID, req.Email, req.Plan)

	return &CheckoutResponse{
		SessionID:   sess.ID,
		CheckoutURL: sess.URL,
	}, nil
}

func (s *StripeService) buildLineItem(plan string) *stripe.CheckoutSessionLineItemParams {
	priceID := s.monthlyPriceID
	if plan == "annual" {
		priceID = s.annualPriceID
	}

	// Si hay Price IDs configurados en Stripe Dashboard, usarlos directamente
	if priceID != "" {
		return &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(priceID),
			Quantity: stripe.Int64(1),
		}
	}

	// Fallback: crear precio ad-hoc con price_data
	unitAmount := int64(1500) // $15.00 USD mensual
	intervalStr := string(stripe.PriceRecurringIntervalMonth)

	if plan == "annual" {
		unitAmount = 15000 // $150.00 USD anual
		intervalStr = string(stripe.PriceRecurringIntervalYear)
	}

	return &stripe.CheckoutSessionLineItemParams{
		PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
			Currency:   stripe.String("usd"),
			UnitAmount: stripe.Int64(unitAmount),
			ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
				Name:        stripe.String("NexoAI Pro — " + plan),
				Description: stripe.String("Acceso completo a NexoAI con generaciones ilimitadas"),
			},
			Recurring: &stripe.CheckoutSessionLineItemPriceDataRecurringParams{
				Interval: stripe.String(intervalStr),
			},
		},
		Quantity: stripe.Int64(1),
	}
}

// ConstructEvent verifica la firma del webhook y devuelve el evento parseado.
func (s *StripeService) ConstructEvent(payload []byte, sigHeader string) (stripe.Event, error) {
	if s.webhookSec == "" {
		return stripe.Event{}, fmt.Errorf("STRIPE_WEBHOOK_SECRET no está configurado")
	}

	event, err := webhook.ConstructEventWithOptions(payload, sigHeader, s.webhookSec,
		webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true},
	)
	if err != nil {
		log.Printf("[Stripe Webhook] Firma inválida: %v", err)
		return stripe.Event{}, fmt.Errorf("firma de webhook inválida: %w", err)
	}

	log.Printf("[Stripe Webhook] Evento verificado: %s", event.Type)
	return event, nil
}
