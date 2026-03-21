package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nexoai/backend/services"
)

type TruoraHandler struct {
	truora *services.TruoraService
}

func NewTruoraHandler(truora *services.TruoraService) *TruoraHandler {
	return &TruoraHandler{truora: truora}
}

func (h *TruoraHandler) GenerateToken(c *gin.Context) {
	var req services.TruoraTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req = services.TruoraTokenRequest{}
	}

	result, err := h.truora.GenerateWebToken(req)
	if err != nil {
		log.Printf("[Truora Handler] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "No se pudo generar el token de verificación",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
