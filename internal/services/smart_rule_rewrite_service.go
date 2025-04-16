package services

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/renatoroquejani/poc-gocache/internal/models"
	"github.com/renatoroquejani/poc-gocache/pkg/gocache"
)

// SmartRuleRewriteService gerencia as Smart Rules de redirecionamento
type SmartRuleRewriteService struct {
	client *gocache.Client
}

// NewSmartRuleRewriteService cria uma nova instu00e2ncia do serviu00e7o de Smart Rules de redirecionamento
func NewSmartRuleRewriteService(client *gocache.Client) *SmartRuleRewriteService {
	return &SmartRuleRewriteService{
		client: client,
	}
}

// GetClient retorna o cliente GoCache usado pelo serviu00e7o
func (s *SmartRuleRewriteService) GetClient() *gocache.Client {
	return s.client
}

// extractURLFromMarkdown extrai a URL real de uma string com formatau00e7u00e3o Markdown
func extractURLFromMarkdown(input string) string {
	// Se nu00e3o tiver formatau00e7u00e3o Markdown, retorna a string original
	if !strings.Contains(input, "[") && !strings.Contains(input, "]") {
		return input
	}

	// Regex para extrair URL de formatos como [url](url) ou [text](url)
	re := regexp.MustCompile(`\[([^\]]*)\]\(([^\)]*)\)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) >= 3 {
		// O segundo grupo de captura (matches[2]) contu00e9m a URL real
		return matches[2]
	}

	// Tenta outro mu00e9todo se o regex nu00e3o funcionar
	if strings.Contains(input, "]") && strings.Contains(input, "(") {
		parts := strings.Split(input, "](")
		if len(parts) > 1 {
			url := parts[1]
			if strings.HasSuffix(url, ")") {
				url = url[:len(url)-1]
			}
			return url
		}
	}

	// Retorna a entrada original se nenhum mu00e9todo funcionar
	return input
}

// CreateRewriteRule cria uma nova regra de redirecionamento
func (s *SmartRuleRewriteService) CreateRewriteRule(request *models.SmartRuleRewriteCreateRequest) (*models.SmartRuleRewriteCreateResponse, error) {
	log.Printf("Criando regra de redirecionamento para domu00ednio %s: %s -> %s",
		request.Domain, request.Match.Request, request.Action.RedirectTo)

	// Constru00f3i os paru00e2metros da requisiu00e7u00e3o
	formData := make(map[string]string)

	// Adiciona os paru00e2metros de match
	if request.Match.RequestURI != "" {
		formData["match[request_uri]"] = request.Match.RequestURI
	} else if request.Match.Request != "" {
		// Para compatibilidade, usa o campo Request se RequestURI nu00e3o for fornecido
		formData["match[request_uri]"] = request.Match.Request
	}

	// Adiciona os mu00e9todos HTTP
	if len(request.Match.RequestMethods) > 0 {
		// Marca o checkbox do mu00e9todo HTTP
		for _, method := range request.Match.RequestMethods {
			formData["match[request_method][]"] = method
		}
	}

	// Adiciona os tipos de dispositivo
	for i, deviceType := range request.Match.DeviceTypes {
		formData[fmt.Sprintf("match[device_type][%d]", i)] = deviceType
	}

	// Adiciona os paru00e2metros de action
	if request.Action.RedirectType != "" {
		formData["action[redirect_type]"] = request.Action.RedirectType
	}

	if request.Action.RedirectTo != "" {
		formData["action[redirect_to]"] = request.Action.RedirectTo
	}

	// Adiciona os novos paru00e2metros de action com os nomes corretos esperados pela API
	if request.Action.RewriteURI != "" {
		formData["action[set_uri]"] = request.Action.RewriteURI
	}

	if request.Action.RewriteHost != "" {
		formData["action[set_host]"] = request.Action.RewriteHost
	}

	if request.Action.Destination != "" {
		formData["action[backend]"] = request.Action.Destination
	}

	// O campo cross_origin na API u00e9 cors
	if request.Action.CrossOrigin != "" {
		// Tratar o problema de formatau00e7u00e3o Markdown
		corsValue := extractURLFromMarkdown(request.Action.CrossOrigin)
		formData["action[cors]"] = corsValue
		// Imprimir para debug
		log.Printf("CORS original: %s, CORS limpo: %s", request.Action.CrossOrigin, corsValue)
	}

	// Configuração SSL completa conforme documentação
	formData["action[ssl_mode]"] = "partial" // Valores válidos: off, flexible, full
	log.Printf("Configurando SSL full com certificado compartilhado")

	// Constru00f3i a URL da requisiu00e7u00e3o
	// Formata o endpoint conforme documentau00e7u00e3o da GoCache
	url := fmt.Sprintf("/rules/settings/%s", request.Domain)

	// Prepara o objeto de resposta
	var response models.SmartRuleRewriteCreateResponse

	// DEBUG: Imprime todos os paru00e2metros da requisiu00e7u00e3o
	log.Printf("Enviando paru00e2metros: %v", formData)

	// Faz a requisiu00e7u00e3o para a API
	resp, err := s.client.Post(url, formData, &response)
	if err != nil {
		log.Printf("Erro ao criar regra de redirecionamento: %v", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("Erro ao criar regra de redirecionamento. Cu00f3digo: %d", resp.StatusCode())
		return nil, fmt.Errorf("erro ao criar regra de redirecionamento. Cu00f3digo: %d", resp.StatusCode())
	}

	log.Printf("Regra de redirecionamento criada com sucesso. ID: %s", response.Response.ID)
	return &response, nil
}

// ListRewriteRules lista todas as regras de redirecionamento de um domu00ednio
func (s *SmartRuleRewriteService) ListRewriteRules(domain string) (*models.SmartRuleRewriteListResponse, error) {
	log.Printf("Listando regras de redirecionamento para domu00ednio %s", domain)

	// Constru00f3i a URL da requisiu00e7u00e3o
	// Formata o endpoint conforme documentau00e7u00e3o da GoCache
	url := fmt.Sprintf("/rules/settings/%s", domain)

	// Prepara o objeto de resposta
	var response models.SmartRuleRewriteListResponse

	// Faz a requisiu00e7u00e3o para a API
	resp, err := s.client.Get(url, &response)
	if err != nil {
		log.Printf("Erro ao listar regras de redirecionamento: %v", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("Erro ao listar regras de redirecionamento. Cu00f3digo: %d", resp.StatusCode())
		return nil, fmt.Errorf("erro ao listar regras de redirecionamento. Cu00f3digo: %d", resp.StatusCode())
	}

	log.Printf("Regras de redirecionamento listadas com sucesso. Total: %d", len(response.Response.Rules))
	return &response, nil
}

// DeleteRewriteRule remove uma regra de redirecionamento
func (s *SmartRuleRewriteService) DeleteRewriteRule(domain, id string) (*models.SmartRuleRewriteDeleteResponse, error) {
	log.Printf("Removendo regra de redirecionamento %s do domu00ednio %s", id, domain)

	// Constru00f3i a URL da requisiu00e7u00e3o
	// Formata o endpoint conforme documentau00e7u00e3o da GoCache
	url := fmt.Sprintf("/rules/settings/%s/%s", domain, id)

	// Prepara o objeto de resposta
	var response models.SmartRuleRewriteDeleteResponse

	// Faz a requisiu00e7u00e3o para a API
	resp, err := s.client.DeleteSimple(url, &response)
	if err != nil {
		log.Printf("Erro ao remover regra de redirecionamento: %v", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("Erro ao remover regra de redirecionamento. Cu00f3digo: %d", resp.StatusCode())
		return nil, fmt.Errorf("erro ao remover regra de redirecionamento. Cu00f3digo: %d", resp.StatusCode())
	}

	log.Printf("Regra de redirecionamento removida com sucesso")
	return &response, nil
}

// CreateSimplifiedRule cria uma regra de redirecionamento padrão com parâmetros simplificados
func (s *SmartRuleRewriteService) CreateSimplifiedRule(request *models.SmartRuleSimplifiedRequest) (*models.SmartRuleRewriteCreateResponse, error) {
	log.Printf("Criando regra de redirecionamento padrão para subdomínio: %s, bucket: %s, conta: %s",
		request.Domain, request.BucketURL, request.AccountID)

	// Valida e ajusta o domínio para ter http:// (requerido pelo CORS da GoCache)
	cors := request.Domain
	if !strings.HasPrefix(cors, "http://") && !strings.HasPrefix(cors, "https://") {
		cors = "http://" + cors // Forçando HTTP como padrão
	}

	// Usa o domínio principal fornecido pelo usuário
	log.Printf("Usando domínio principal %s para criar regra no subdomínio %s", request.ParentDomain, request.Domain)

	// Cria o request completo para a API
	completeRequest := &models.SmartRuleRewriteCreateRequest{
		// Usa o domínio principal como alvo, mas configura o host para o subdomínio
		Domain: request.ParentDomain,
		Match: models.SmartRuleRewriteMatch{
			RequestURI: "/*",
			// Adiciona o host para especificar o subdomínio
			Host: request.Domain,
		},
		Action: models.SmartRuleRewriteAction{
			CrossOrigin: cors,
			RewriteURI:  "/" + request.AccountID + "/$1",
			RewriteHost: request.BucketURL,
			Destination: request.BucketURL,
		},
	}

	// Log detalhado da configuração que será enviada para GoCache
	log.Printf("DETALHES DA CONFIGURAÇÃO GOCACHE: Enviando para domínio principal '%s', com Host='%s', RewriteHost='%s', BucketURL='%s'",
		completeRequest.Domain,
		completeRequest.Match.Host,
		completeRequest.Action.RewriteHost,
		completeRequest.Action.Destination)

	// Usa o método existente para criar a regra
	return s.CreateRewriteRule(completeRequest)
}

// UpdateRewriteRule atualiza uma regra de redirecionamento
func (s *SmartRuleRewriteService) UpdateRewriteRule(domain, id string, request *models.SmartRuleRewriteCreateRequest) (*models.SmartRuleRewriteUpdateResponse, error) {
	log.Printf("Atualizando regra de redirecionamento %s do domu00ednio %s", id, domain)

	// Constru00f3i os paru00e2metros da requisiu00e7u00e3o
	formData := make(map[string]string)

	// Adiciona os paru00e2metros de match
	if request.Match.RequestURI != "" {
		formData["match[request_uri]"] = request.Match.RequestURI
	} else if request.Match.Request != "" {
		// Para compatibilidade, usa o campo Request se RequestURI nu00e3o for fornecido
		formData["match[request_uri]"] = request.Match.Request
	}

	// Adiciona os mu00e9todos HTTP
	if len(request.Match.RequestMethods) > 0 {
		// Marca o checkbox do mu00e9todo HTTP
		for _, method := range request.Match.RequestMethods {
			formData["match[request_method][]"] = method
		}
	}

	// Adiciona os tipos de dispositivo
	for i, deviceType := range request.Match.DeviceTypes {
		formData[fmt.Sprintf("match[device_type][%d]", i)] = deviceType
	}

	// Adiciona os paru00e2metros de action
	if request.Action.RedirectType != "" {
		formData["action[redirect_type]"] = request.Action.RedirectType
	}

	if request.Action.RedirectTo != "" {
		formData["action[redirect_to]"] = request.Action.RedirectTo
	}

	// Adiciona os novos paru00e2metros de action com os nomes corretos esperados pela API
	if request.Action.RewriteURI != "" {
		formData["action[set_uri]"] = request.Action.RewriteURI
	}

	if request.Action.RewriteHost != "" {
		formData["action[set_host]"] = request.Action.RewriteHost
	}

	if request.Action.Destination != "" {
		formData["action[backend]"] = request.Action.Destination
	}

	// O campo cross_origin na API u00e9 cors
	if request.Action.CrossOrigin != "" {
		// Tratar o problema de formatau00e7u00e3o Markdown
		corsValue := extractURLFromMarkdown(request.Action.CrossOrigin)
		formData["action[cors]"] = corsValue
	}

	// Constru00f3i a URL da requisiu00e7u00e3o
	// Formata o endpoint conforme documentau00e7u00e3o da GoCache
	url := fmt.Sprintf("/rules/settings/%s/%s", domain, id)

	// Prepara o objeto de resposta
	var response models.SmartRuleRewriteUpdateResponse

	// Faz a requisiu00e7u00e3o para a API
	resp, err := s.client.Put(url, formData, &response)
	if err != nil {
		log.Printf("Erro ao atualizar regra de redirecionamento: %v", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("Erro ao atualizar regra de redirecionamento. Cu00f3digo: %d", resp.StatusCode())
		return nil, fmt.Errorf("erro ao atualizar regra de redirecionamento. Cu00f3digo: %d", resp.StatusCode())
	}

	log.Printf("Regra de redirecionamento atualizada com sucesso")
	return &response, nil
}
