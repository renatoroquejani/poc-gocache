package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/renatoroquejani/poc-gocache/internal/models"
	"github.com/renatoroquejani/poc-gocache/internal/services"
)

// DomainHandler handles domain + smart rule operations
type DomainHandler struct {
	domainService    *services.DomainService
	smartRuleService *services.SmartRuleService
}

// NewDomainHandler creates a new DomainHandler
func NewDomainHandler(domainService *services.DomainService, smartRuleService *services.SmartRuleService) *DomainHandler {
	return &DomainHandler{
		domainService:    domainService,
		smartRuleService: smartRuleService,
	}
}

// RegisterRoutes registers domain routes
func (h *DomainHandler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/domains")
	{
		group.GET("", h.ListDomains)
		group.POST("", h.CreateDomain)
		group.DELETE("/:domainID", h.DeleteDomainWithSmartRules)
	}

	rulesGroup := router.Group("/rules")
	{
		rulesGroup.POST("", h.CreateSmartRule)
	}
}

// CreateDomainWithSmartRule godoc
// @Summary Create domain with smart rule
// @Description Creates a domain and an associated smart rule
// @Tags Domains
// @Accept json
// @Produce json
// @Param request body models.DomainCreateRequest true "Domain info"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /domains [post]
func (h *DomainHandler) CreateDomain(c *gin.Context) {
	var req models.DomainCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.domainService.CreateDomain(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *DomainHandler) CreateSmartRule(c *gin.Context) {
	// Serviço foi removido/substituído pelo SmartRuleRewriteService
	c.JSON(http.StatusNotFound, gin.H{"error": "Este endpoint foi descontinuado. Use /rules/settings/{domain} em vez disso."})
}

// ListDomains godoc
// @Summary Listar domínios
// @Description Lista todos os domínios disponíveis na GoCache
// @Tags Domains
// @Produce json
// @Success 200 {object} models.DomainListResponse
// @Failure 500 {object} map[string]interface{}
// @Router /domains [get]
func (h *DomainHandler) ListDomains(c *gin.Context) {
	domains, err := h.domainService.ListDomains()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao listar domínios: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, domains)
}

// DeleteDomainWithSmartRules godoc
// @Summary Delete domain
// @Description Deletes a domain
// @Tags Domains
// @Param domainID path int true "Domain ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /domains/{domainID} [delete]
func (h *DomainHandler) DeleteDomainWithSmartRules(c *gin.Context) {
	domainIDStr := c.Param("domainID")
	domainID, err := strconv.Atoi(domainIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid domain ID"})
		return
	}

	// A funcionalidade de listar e excluir Smart Rules foi movida para outro endpoint
	// Apenas excluir o domínio
	err = h.domainService.DeleteDomain(domainID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete domain: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "domain deleted"})
}
