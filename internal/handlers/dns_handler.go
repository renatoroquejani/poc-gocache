package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/renatoroquejani/poc-gocache/internal/models"
	"github.com/renatoroquejani/poc-gocache/internal/services"
)

// DNSHandler manipula as requisições relacionadas a domínios
type DNSHandler struct {
	service *services.DNSService
}

// NewDNSHandler cria uma nova instância de DNSHandler
func NewDNSHandler(service *services.DNSService) *DNSHandler {
	return &DNSHandler{
		service: service,
	}
}

// RegisterRoutes registra as rotas no router do Gin
func (h *DNSHandler) RegisterRoutes(router *gin.RouterGroup) {
	dnsGroup := router.Group("/dns")
	{
		dnsGroup.GET("", h.ListDNS)
		dnsGroup.GET("/:id", h.GetDNS)
		dnsGroup.POST("/:domain", h.CreateDNS)
		dnsGroup.PUT("/:id", h.UpdateDNS)
		dnsGroup.DELETE("/:id", h.DeleteDNS)
	}
}

// ListDNS godoc
// @Summary Lista todos os registros DNS
// @Description Retorna uma lista de todos os registros DNS cadastrados para um domínio específico
// @Tags DNS
// @Accept json
// @Produce json
// @Param domain query string true "Domínio para listar os registros DNS"
// @Success 200 {object} models.DNSListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /dns [get]
func (h *DNSHandler) ListDNS(c *gin.Context) {
	domain := c.Query("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Domínio não especificado"})
		return
	}

	response, err := h.service.ListDNS(domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetDNS godoc
// @Summary Obtém um registro DNS específico
// @Description Retorna os detalhes de um registro DNS específico
// @Tags DNS
// @Accept json
// @Produce json
// @Param id path int true "ID do registro DNS"
// @Success 200 {object} models.DNSCreateResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /dns/{id} [get]
func (h *DNSHandler) GetDNS(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	response, err := h.service.GetDNS(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateDNS godoc
// @Summary Cria um novo registro DNS
// @Description Cria um novo registro DNS na Gocache
// @Tags DNS
// @Accept json
// @Produce json
// @Param domain path string true "Domínio para o qual criar o registro DNS"
// @Param request body models.DNSCreateRequest true "Dados do registro DNS"
// @Success 201 {object} models.DNSCreateResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /dns/{domain} [post]
func (h *DNSHandler) CreateDNS(c *gin.Context) {
	var request models.DNSCreateRequest
	contentType := c.GetHeader("Content-Type")

	if contentType == "application/x-www-form-urlencoded" || contentType == "application/x-www-form-urlencoded; charset=UTF-8" {
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// Extract domain from URL
	domain := c.Param("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Domínio não especificado"})
		return
	}
	request.Domain = domain

	// Only create DNS record (assumes domain already exists in GoCache)
	response, err := h.service.CreateDNS(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// UpdateDNS godoc
// @Summary Atualiza um registro DNS existente
// @Description Atualiza os dados de um registro DNS existente na Gocache
// @Tags DNS
// @Accept json
// @Produce json
// @Param id path int true "ID do registro DNS"
// @Param request body models.DNSUpdateRequest true "Dados do domínio"
// @Success 200 {object} models.DNSUpdateResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /dns/{id} [put]
func (h *DNSHandler) UpdateDNS(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var request models.DNSUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.UpdateDNS(id, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteDNS godoc
// @Summary Remove um registro DNS
// @Description Remove um registro DNS existente na Gocache
// @Tags DNS
// @Accept json
// @Produce json
// @Param id path int true "ID do registro DNS"
// @Success 200 {object} models.DNSDeleteResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /dns/{id} [delete]
func (h *DNSHandler) DeleteDNS(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	response, err := h.service.DeleteDNS(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
