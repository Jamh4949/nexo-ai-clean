package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/nexoai/backend/config"
)

// Endpoint correcto para generar Web Integration Tokens (JWT de 2h)
// Documentación: https://dev.truora.com/guides/web_integration_token/
const truoraTokenURL = "https://api.account.truora.com/v1/api-keys"

const truoraIframeBaseURL = "https://identity.truora.com/"

type TruoraService struct {
	apiKey string
	flowID string
}

func NewTruoraService(cfg *config.Config) *TruoraService {
	return &TruoraService{
		apiKey: cfg.TruoraAPIKey,
		flowID: cfg.TruoraFlowID,
	}
}

type TruoraTokenRequest struct {
	AccountID string `json:"account_id"`
}

type TruoraTokenResponse struct {
	Success    bool   `json:"success"`
	ProcessURL string `json:"process_url"`
}

func (s *TruoraService) GenerateWebToken(req TruoraTokenRequest) (*TruoraTokenResponse, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("TRUORA_API_KEY no está configurada")
	}
	if s.flowID == "" {
		return nil, fmt.Errorf("TRUORA_FLOW_ID no está configurado")
	}

	accountID := req.AccountID
	if accountID == "" {
		accountID = uuid.NewString()
	}

	form := url.Values{}
	form.Set("key_type", "web")
	form.Set("grant", "digital-identity")
	form.Set("api_key_version", "1")
	form.Set("country", "ALL")
	form.Set("flow_id", s.flowID)
	form.Set("account_id", accountID)

	httpReq, err := http.NewRequest("POST", truoraTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creando request de Truora: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Truora-API-Key", s.apiKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error conectando con Truora: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta de Truora: %w", err)
	}

	if resp.StatusCode >= 400 {
		log.Printf("[Truora] Error %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("Truora respondió con status %d", resp.StatusCode)
	}

	var result struct {
		APIKey  string `json:"api_key"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("[Truora] Respuesta no parseable: %s", string(body))
		return nil, fmt.Errorf("error parseando respuesta de Truora: %w", err)
	}

	if result.APIKey == "" {
		log.Printf("[Truora] Respuesta sin api_key: %s", string(body))
		return nil, fmt.Errorf("Truora no devolvió un token válido")
	}

	processURL := fmt.Sprintf("%s?token=%s", truoraIframeBaseURL, result.APIKey)

	log.Printf("[Truora] Web Integration Token generado para account: %s", accountID)

	return &TruoraTokenResponse{
		Success:    true,
		ProcessURL: processURL,
	}, nil
}
