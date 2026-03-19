package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/nexoai/backend/config"
)

const hubspotContactsURL = "https://api.hubapi.com/crm/v3/objects/contacts"

type CRMService struct {
	apiKey string
}

func NewCRMService(cfg *config.Config) *CRMService {
	return &CRMService{apiKey: cfg.HubSpotAPIKey}
}

type ContactRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

func (s *CRMService) CreateContact(req ContactRequest) error {
	if s.apiKey == "" {
		log.Printf("[CRM] HUBSPOT_API_KEY no configurada — contacto omitido para %s", req.Email)
		return nil
	}

	payload := map[string]interface{}{
		"properties": map[string]string{
			"email":     req.Email,
			"firstname": req.FirstName,
			"lastname":  req.LastName,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error serializando payload de HubSpot: %w", err)
	}

	httpReq, err := http.NewRequest("POST", hubspotContactsURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creando request de HubSpot: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error conectando con HubSpot: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		log.Printf("[CRM] Contacto ya existe en HubSpot: %s", req.Email)
		return nil
	}

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		log.Printf("[CRM] HubSpot error %d: %s", resp.StatusCode, string(respBody))
		return fmt.Errorf("HubSpot respondió con status %d: %s", resp.StatusCode, string(respBody))
	}

	log.Printf("[CRM] Contacto creado en HubSpot: %s", req.Email)
	return nil
}
