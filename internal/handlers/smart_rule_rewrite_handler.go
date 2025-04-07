package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/renatoroquejani/poc-gocache/internal/models"
	"github.com/renatoroquejani/poc-gocache/internal/services"
)

// SmartRuleRewriteHandler gerencia as requisiu00e7u00f5es relacionadas u00e0s Smart Rules de redirecionamento
type SmartRuleRewriteHandler struct {
	service *services.SmartRuleRewriteService
}

// NewSmartRuleRewriteHandler cria uma nova instu00e2ncia do handler de Smart Rules de redirecionamento
func NewSmartRuleRewriteHandler(service *services.SmartRuleRewriteService) *SmartRuleRewriteHandler {
	return &SmartRuleRewriteHandler{
		service: service,
	}
}

// CreateRewriteRule cria uma nova regra de redirecionamento
// @Summary Criar uma nova regra de redirecionamento
// @Description Cria uma nova regra de redirecionamento para um domínio específico
// @Tags Smart Rules
// @Accept json
// @Produce json
// @Param domain path string true "Nome do domínio"
// @Param request body models.SmartRuleRewriteCreateRequest true "Dados da regra de redirecionamento"
// @Success 200 {object} models.SmartRuleRewriteCreateResponse
// @Failure 400 {object} map[string]interface{} "Erro na requisição"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /rules/settings/{domain} [post]
func (h *SmartRuleRewriteHandler) CreateRewriteRule(c *gin.Context) {
	var request models.SmartRuleRewriteCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtu00e9m o domu00ednio da URL
	domain := c.Param("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domu00ednio nu00e3o especificado"})
		return
	}

	// Define o domu00ednio na requisiu00e7u00e3o
	request.Domain = domain

	// Cria a regra de redirecionamento
	response, err := h.service.CreateRewriteRule(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListRewriteRules lista todas as regras de redirecionamento de um domu00ednio
// @Summary Listar regras de redirecionamento
// @Description Lista todas as regras de redirecionamento para um domínio específico
// @Tags Smart Rules
// @Accept json
// @Produce json
// @Param domain path string true "Nome do domínio"
// @Success 200 {object} models.SmartRuleRewriteListResponse
// @Failure 400 {object} map[string]interface{} "Erro na requisição"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /rules/settings/{domain} [get]
func (h *SmartRuleRewriteHandler) ListRewriteRules(c *gin.Context) {
	// Obtu00e9m o domu00ednio da URL
	domain := c.Param("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domu00ednio nu00e3o especificado"})
		return
	}

	// Lista as regras de redirecionamento
	response, err := h.service.ListRewriteRules(domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteRewriteRule remove uma regra de redirecionamento
// @Summary Remover uma regra de redirecionamento
// @Description Remove uma regra de redirecionamento específica de um domínio
// @Tags Smart Rules
// @Accept json
// @Produce json
// @Param domain path string true "Nome do domínio"
// @Param id path string true "ID da regra de redirecionamento"
// @Success 200 {object} models.SmartRuleRewriteDeleteResponse
// @Failure 400 {object} map[string]interface{} "Erro na requisição"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /rules/settings/{domain}/{id} [delete]
func (h *SmartRuleRewriteHandler) DeleteRewriteRule(c *gin.Context) {
	// Obtu00e9m o domu00ednio e o ID da regra da URL
	domain := c.Param("domain")
	id := c.Param("id")

	if domain == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domu00ednio ou ID nu00e3o especificado"})
		return
	}

	// Remove a regra de redirecionamento
	response, err := h.service.DeleteRewriteRule(domain, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateRewriteRule atualiza uma regra de redirecionamento
// @Summary Atualizar uma regra de redirecionamento
// @Description Atualiza uma regra de redirecionamento específica de um domínio
// @Tags Smart Rules
// @Accept json
// @Produce json
// @Param domain path string true "Nome do domínio"
// @Param id path string true "ID da regra de redirecionamento"
// @Param request body models.SmartRuleRewriteCreateRequest true "Dados da regra de redirecionamento"
// @Success 200 {object} models.SmartRuleRewriteUpdateResponse
// @Failure 400 {object} map[string]interface{} "Erro na requisição"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /rules/settings/{domain}/{id} [put]
func (h *SmartRuleRewriteHandler) UpdateRewriteRule(c *gin.Context) {
	var request models.SmartRuleRewriteCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtu00e9m o domu00ednio e o ID da regra da URL
	domain := c.Param("domain")
	id := c.Param("id")

	if domain == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domu00ednio ou ID nu00e3o especificado"})
		return
	}

	// Define o domu00ednio na requisiu00e7u00e3o
	request.Domain = domain

	// Atualiza a regra de redirecionamento
	response, err := h.service.UpdateRewriteRule(domain, id, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RegisterRoutes registra as rotas do handler no router
// CreateSimplifiedRule cria uma regra padrão de redirecionamento com parâmetros simplificados
// @Summary Criar regra padrão de redirecionamento
// @Description Cria uma regra padrão de redirecionamento com parâmetros simplificados
// @Tags Smart Rules
// @Accept json
// @Produce json
// @Param domain path string false "Domínio principal (ex: exod.com.br)"
// @Param request body models.SmartRuleSimplifiedRequest true "Parâmetros simplificados"
// @Success 200 {object} models.SmartRuleRewriteCreateResponse
// @Failure 400 {object} map[string]interface{} "Erro na requisição"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /rules/{domain}/simplified [post]
func (h *SmartRuleRewriteHandler) CreateSimplifiedRule(c *gin.Context) {
	// Verifica se o domínio foi fornecido na URL
	domain := c.Param("domain")
	
	var request models.SmartRuleSimplifiedRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Se o domínio foi fornecido na URL, usa-o como parent_domain
	if domain != "" {
		request.ParentDomain = domain
	}

	// Valida se temos o parent_domain (ou da URL ou do body)
	if request.ParentDomain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parent_domain não foi especificado nem na URL nem no corpo da requisição"})
		return
	}

	// Cria a regra de redirecionamento simplificada
	response, err := h.service.CreateSimplifiedRule(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetSimplifiedRuleForm godoc
// @Summary Obter formulário para criação de regra simplificada
// @Description Retorna os dados necessários para criar uma regra simplificada, incluindo a lista de domínios disponíveis
// @Tags Smart Rules
// @Produce json
// @Success 200 {object} models.SmartRuleSimplifiedFormResponse
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /rules/simplified/form [get]

// CreateSimplifiedRuleWithDomain godoc
// @Summary Criar regra padrão de redirecionamento usando domínio da URL
// @Description Cria uma regra padrão de redirecionamento com domínio especificado na URL e parâmetros simplificados no body
// @Tags Smart Rules
// @Accept json
// @Produce json
// @Param domain path string true "Domínio principal (ex: exod.com.br)"
// @Param request body models.SmartRuleSimplifiedRequest true "Parâmetros simplificados (sem parent_domain)"
// @Success 200 {object} models.SmartRuleRewriteCreateResponse
// @Failure 400 {object} map[string]interface{} "Erro na requisição"
// @Failure 500 {object} map[string]interface{} "Erro interno do servidor"
// @Router /rules/{domain}/simplified [post]
func (h *SmartRuleRewriteHandler) CreateSimplifiedRuleWithDomain(c *gin.Context) {
	// Obtém o domínio da URL
	domain := c.Param("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domínio não especificado na URL"})
		return
	}

	// Obtém o corpo da requisição
	var request models.SmartRuleSimplifiedRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Define o domínio principal com o valor da URL
	request.ParentDomain = domain

	// Cria a regra de redirecionamento simplificada
	response, err := h.service.CreateSimplifiedRule(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *SmartRuleRewriteHandler) GetSimplifiedRuleForm(c *gin.Context) {
	// Precisamos obter a lista de domínios através do DomainService
	// Como não temos acesso direto, podemos usar o client do SmartRuleService
	// para fazer a chamada diretamente
	
	// Faz a requisição para API de domínios da GoCache
	var domainResponse models.DomainListResponse
	endpoint := "/domain"
	
	_, err := h.service.GetClient().Get(endpoint, &domainResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao obter domínios: " + err.Error()})
		return
	}
	
	// Converte para o formato esperado na resposta
	domainOptions := make([]models.DomainOption, 0)
	for _, domainName := range domainResponse.Response.Domains {
		domainOptions = append(domainOptions, models.DomainOption{
			Name:        domainName,
			DisplayName: domainName,
		})
	}
	
	// Retorna a resposta
	c.JSON(http.StatusOK, models.SmartRuleSimplifiedFormResponse{
		Domains: domainOptions,
	})
}

func (h *SmartRuleRewriteHandler) RegisterRoutes(router gin.IRouter) {
	group := router.Group("/rules/settings")
	{
		group.POST("/:domain", h.CreateRewriteRule)
		group.GET("/:domain", h.ListRewriteRules)
		group.DELETE("/:domain/:id", h.DeleteRewriteRule)
		group.PUT("/:domain/:id", h.UpdateRewriteRule)
	}

	// Rotas para criação simplificada de regras
	router.POST("/rules/:domain/simplified", h.CreateSimplifiedRule)
	router.GET("/rules/simplified/form", h.GetSimplifiedRuleForm)
}
