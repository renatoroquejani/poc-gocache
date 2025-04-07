package handlers

import (
	"net/http"
	"strconv"

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
		cacheGroup.POST("/purge", h.PurgeCache)
		cacheGroup.POST("/purge-by-prefix", h.PurgeCacheByPrefix)
		cacheGroup.GET("/status/:domainId", h.GetCacheStatus)
	}
}

// PurgeCache godoc
// @Summary Expira o cache de URLs específicas
// @Description Remove o cache de URLs específicas para um domínio
// @Tags Cache
// @Accept json
// @Produce json
// @Param request body models.CachePurgeRequest true "Dados para expiração de cache"
// @Success 200 {object} models.CacheInvalidationResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cache/purge [post]
func (h *CacheHandler) PurgeCache(c *gin.Context) {
	var request models.CachePurgeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(request.URLs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A lista de URLs não pode estar vazia"})
		return
	}

	response, err := h.service.PurgeCache(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// PurgeCacheByPrefix godoc
// @Summary Expira o cache por prefixo de URL
// @Description Remove o cache de todas as URLs que começam com um prefixo específico
// @Tags Cache
// @Accept json
// @Produce json
// @Param request body models.CachePurgeByPrefixRequest true "Dados para expiração de cache por prefixo"
// @Success 200 {object} models.CacheInvalidationResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cache/purge-by-prefix [post]
func (h *CacheHandler) PurgeCacheByPrefix(c *gin.Context) {
	var request models.CachePurgeByPrefixRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Prefix == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "O prefixo não pode estar vazio"})
		return
	}

	response, err := h.service.PurgeCacheByPrefix(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetCacheStatus godoc
// @Summary Obtém o status do cache para um domínio
// @Description Retorna informações sobre o status do cache para um domínio específico
// @Tags Cache
// @Accept json
// @Produce json
// @Param domainId path int true "ID do domínio"
// @Success 200 {object} models.CacheStatusResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /cache/status/{domainId} [get]
func (h *CacheHandler) GetCacheStatus(c *gin.Context) {
	domainIdStr := c.Param("domainId")
	domainId, err := strconv.Atoi(domainIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de domínio inválido"})
		return
	}

	response, err := h.service.GetCacheStatus(domainId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
