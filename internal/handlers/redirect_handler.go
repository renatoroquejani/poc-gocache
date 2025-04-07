package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/renatoroquejani/poc-gocache/internal/models"
	"github.com/renatoroquejani/poc-gocache/internal/services"
)

// RedirectHandler gerencia as requisições relacionadas a regras de redirecionamento
type RedirectHandler struct {
	service *services.RedirectService
}

// NewRedirectHandler cria uma nova instância do handler de redirecionamento
func NewRedirectHandler(service *services.RedirectService) *RedirectHandler {
	return &RedirectHandler{
		service: service,
	}
}

// CreateRedirect cria uma nova regra de redirecionamento
func (h *RedirectHandler) CreateRedirect(c *gin.Context) {
	var request models.RedirectCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.CreateRedirect(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// ListRedirects lista todas as regras de redirecionamento para um domínio
func (h *RedirectHandler) ListRedirects(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domínio não especificado"})
		return
	}

	response, err := h.service.ListRedirects(domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteRedirect exclui uma regra de redirecionamento
func (h *RedirectHandler) DeleteRedirect(c *gin.Context) {
	domain := c.Param("domain")
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	response, err := h.service.DeleteRedirect(domain, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RegisterRoutes registra as rotas do handler no router
func (h *RedirectHandler) RegisterRoutes(router *gin.Engine) {
	redirectGroup := router.Group("/api/v1/redirects")
	{
		redirectGroup.POST("", h.CreateRedirect)
		redirectGroup.GET("", h.ListRedirects)
		redirectGroup.DELETE("/:domain/:id", h.DeleteRedirect)
	}
}
