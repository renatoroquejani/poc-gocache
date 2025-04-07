package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/renatoroquejani/poc-gocache/internal/models"
	"github.com/renatoroquejani/poc-gocache/internal/services"
)

// CacheHandler manipula as requisições relacionadas a cache
type CacheHandler struct {
	service *services.CacheService
}

// NewCacheHandler cria uma nova instância de CacheHandler
func NewCacheHandler(service *services.CacheService) *CacheHandler {
	return &CacheHandler{
		service: service,
	}
}

// RegisterRoutes registra as rotas no router do Gin
func (h *CacheHandler) RegisterRoutes(router *gin.RouterGroup) {
	cacheGroup := router.Group("/cache")
	{
		cacheGroup.DELETE("/purge-all/:domainName", h.PurgeAllCache)
		cacheGroup.DELETE("/purge-urls", h.PurgeUrls)
	}
}

// PurgeUrls godoc
// @Summary Expira o cache de URLs específicas
// @Description Remove o cache de URLs específicas para um domínio, podendo incluir wildcards
// @Tags Cache
// @Accept json
// @Produce json
// @Param request body models.CachePurgeRequest true "Dados para expiração de cache"
// @Success 200 {object} models.CacheInvalidationResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cache/purge-urls [delete]
func (h *CacheHandler) PurgeUrls(c *gin.Context) {
	var request models.CachePurgeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(request.URLs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A lista de URLs não pode estar vazia"})
		return
	}

	response, err := h.service.PurgeUrls(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// PurgeAllCache godoc
// @Summary Expira todo o cache de um domínio
// @Description Remove todo o cache de um domínio específico
// @Tags Cache
// @Accept json
// @Produce json
// @Param domainName path string true "Nome do domínio"
// @Success 200 {object} models.CacheInvalidationResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cache/purge-all/{domainName} [delete]
func (h *CacheHandler) PurgeAllCache(c *gin.Context) {
	domainName := c.Param("domainName")
	if domainName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome de domínio inválido"})
		return
	}

	response, err := h.service.PurgeAllCache(domainName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}


