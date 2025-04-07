package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/renatoroquejani/poc-gocache/internal/models"
	"github.com/renatoroquejani/poc-gocache/internal/services"
)

// ProxyHandler gerencia as requisiu00e7u00f5es relacionadas ao proxy de redirecionamento
type ProxyHandler struct {
	service *services.ProxyService
}

// NewProxyHandler cria uma nova instu00e2ncia do handler de proxy
func NewProxyHandler(service *services.ProxyService) *ProxyHandler {
	return &ProxyHandler{
		service: service,
	}
}

// AddMapping adiciona um novo mapeamento de domu00ednio
func (h *ProxyHandler) AddMapping(c *gin.Context) {
	var mapping models.DomainMapping
	if err := c.ShouldBindJSON(&mapping); err != nil {
		c.JSON(http.StatusBadRequest, models.DomainMappingResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	err := h.service.AddMapping(mapping)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.DomainMappingResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.DomainMappingResponse{
		Success: true,
		Mapping: mapping,
	})
}

// GetMappings lista todos os mapeamentos de domu00ednio
func (h *ProxyHandler) GetMappings(c *gin.Context) {
	mappings := h.service.GetAllMappings()

	c.JSON(http.StatusOK, models.DomainMappingsListResponse{
		Success:  true,
		Mappings: mappings,
		Total:    len(mappings),
	})
}

// DeleteMapping remove um mapeamento de domu00ednio
func (h *ProxyHandler) DeleteMapping(c *gin.Context) {
	domain := c.Param("domain")

	err := h.service.DeleteMapping(domain)
	if err != nil {
		c.JSON(http.StatusNotFound, models.DomainMappingResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.DomainMappingResponse{
		Success: true,
	})
}

// HandleRedirect processa redirecionamentos com base nos mapeamentos configurados
func (h *ProxyHandler) HandleRedirect(c *gin.Context) {
	host := c.Request.Host
	path := c.Request.URL.Path

	log.Printf("Recebida requisiu00e7u00e3o para host: %s, path: %s", host, path)

	// Remove a porta do host, se presente
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	// Procura o mapeamento correspondente
	mapping, err := h.service.GetMapping(host)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "domu00ednio nu00e3o configurado"})
		return
	}

	destination := mapping.Destination

	// Se o path nu00e3o for a raiz, adiciona ao destino
	if path != "/" {
		// Remove a barra inicial do path se o destino ju00e1 terminar com barra
		if strings.HasSuffix(destination, "/") && strings.HasPrefix(path, "/") {
			path = path[1:]
		}
		destination = destination + path
	}

	log.Printf("Redirecionando para: %s", destination)
	c.Redirect(http.StatusMovedPermanently, destination)
}

// RegisterRoutes registra as rotas do handler no router
func (h *ProxyHandler) RegisterRoutes(router *gin.Engine) {
	proxyGroup := router.Group("/api/v1/proxy")
	{
		proxyGroup.POST("/mappings", h.AddMapping)
		proxyGroup.GET("/mappings", h.GetMappings)
		proxyGroup.DELETE("/mappings/:domain", h.DeleteMapping)
	}

	// Rota para processar redirecionamentos (deve ser registrada separadamente no main.go)
}
