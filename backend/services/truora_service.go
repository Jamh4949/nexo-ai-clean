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

const truoraProcessesURL = "https://api.identity.truora.com/v1/processes"

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
	ProcessID  string `json:"process_id"`
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
	form.Set("type", "web")
	form.Set("flow_id", s.flowID)
	form.Set("account_id", accountID)

	httpReq, err := http.NewRequest("POST", truoraProcessesURL, strings.NewReader(form.Encode()))
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

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parseando respuesta de Truora: %w", err)
	}

	processURL, _ := result["web_url"].(string)
	if processURL == "" {
		processURL, _ = result["process_url"].(string)
	}

	processID, _ := result["process_id"].(string)

	if processURL == "" {
		log.Printf("[Truora] Respuesta sin URL de proceso: %s", string(body))
		return nil, fmt.Errorf("Truora no devolvió una URL de proceso válida")
	}

	log.Printf("[Truora] Proceso creado: %s (account: %s)", processID, accountID)

	return &TruoraTokenResponse{
		Success:    true,
		ProcessURL: processURL,
		ProcessID:  processID,
	}, nil
}
