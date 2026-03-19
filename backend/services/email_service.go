package services

import (
	"fmt"
	"log"

	"github.com/nexoai/backend/config"
	"github.com/resend/resend-go/v2"
)

type EmailService struct {
	client *resend.Client
}

func NewEmailService(cfg *config.Config) *EmailService {
	var client *resend.Client
	if cfg.ResendAPIKey != "" {
		client = resend.NewClient(cfg.ResendAPIKey)
	}
	return &EmailService{client: client}
}

type WelcomeEmailRequest struct {
	To   string `json:"to"`
	Name string `json:"name"`
}

func (s *EmailService) SendWelcomeEmail(req WelcomeEmailRequest) error {
	if s.client == nil {
		log.Printf("[Email] RESEND_API_KEY no configurada — email omitido para %s", req.To)
		return nil
	}

	displayName := req.Name
	if displayName == "" {
		displayName = req.To
	}

	html := fmt.Sprintf(`
		<div style="font-family: 'Inter', sans-serif; max-width: 600px; margin: 0 auto; padding: 40px 20px;">
			<div style="text-align: center; margin-bottom: 32px;">
				<div style="display: inline-block; background: #4f46e5; color: white; font-weight: bold; padding: 8px 12px; border-radius: 8px; font-size: 18px;">N</div>
				<span style="font-size: 24px; font-weight: bold; margin-left: 8px;">NexoAI</span>
			</div>
			<h1 style="color: #111827; font-size: 24px;">¡Hola %s!</h1>
			<p style="color: #4b5563; font-size: 16px; line-height: 1.6;">
				Gracias por unirte a <strong>NexoAI</strong>. Estás a un paso de generar contenido
				increíble para tus redes sociales con inteligencia artificial.
			</p>
			<p style="color: #4b5563; font-size: 16px; line-height: 1.6;">
				Tu suscripción está activa y puedes comenzar a crear contenido ahora mismo.
			</p>
			<div style="text-align: center; margin: 32px 0;">
				<a href="https://nexoai.com/dashboard" style="background: #4f46e5; color: white; padding: 12px 32px; border-radius: 8px; text-decoration: none; font-weight: 600;">
					Ir al Dashboard
				</a>
			</div>
			<p style="color: #9ca3af; font-size: 14px;">— El equipo de NexoAI</p>
		</div>`, displayName)

	params := &resend.SendEmailRequest{
		From:    "NexoAI <onboarding@resend.dev>",
		To:      []string{req.To},
		Subject: "¡Bienvenido a NexoAI!",
		Html:    html,
	}

	sent, err := s.client.Emails.Send(params)
	if err != nil {
		log.Printf("[Email] Error enviando a %s: %v", req.To, err)
		return fmt.Errorf("error enviando email via Resend: %w", err)
	}

	log.Printf("[Email] Correo de bienvenida enviado a %s (ID: %s)", req.To, sent.Id)
	return nil
}
