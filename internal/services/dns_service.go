package services

import (
	"fmt"

	"github.com/renatoroquejani/poc-gocache/internal/models"
	"github.com/renatoroquejani/poc-gocache/pkg/gocache"
)

// DNSService fornece métodos para interagir com a API de domínios da Gocache
type DNSService struct {
	Client *gocache.Client
}

// NewDNSService cria uma nova instância de DNSService
func NewDNSService(client *gocache.Client) *DNSService {
	return &DNSService{
		Client: client,
	}
}

// ListDNS lista todos os domínios cadastrados para um domínio específico
func (s *DNSService) ListDNS(domain string) (*models.DNSListResponse, error) {
	if domain == "" {
		return nil, fmt.Errorf("domínio não especificado")
	}

	// Endpoint correto conforme documentação da GoCache
	endpoint := fmt.Sprintf("/dns/%s", domain)
	result := &models.DNSListResponse{}

	_, err := s.Client.Get(endpoint, result)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar domínios: %w", err)
	}

	return result, nil
}

// GetDNS obtém detalhes de um domínio específico pelo ID
func (s *DNSService) GetDNS(id int) (*models.DNSCreateResponse, error) {
	// Endpoint correto conforme documentação da GoCache
	endpoint := fmt.Sprintf("/dns/%d", id)
	result := &models.DNSCreateResponse{}

	_, err := s.Client.Get(endpoint, result)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter domínio: %w", err)
	}

	return result, nil
}

// CreateDNS cria um novo domínio
func (s *DNSService) CreateDNS(req models.DNSCreateRequest) (*models.DNSCreateResponse, error) {
	if req.Domain == "" {
		return nil, fmt.Errorf("domínio não especificado")
	}

	// Formata o endpoint conforme documentação: /v1/dns/{dominio}
	endpoint := fmt.Sprintf("/dns/%s", req.Domain)
	result := &models.DNSCreateResponse{}

	_, err := s.Client.Post(endpoint, req, result)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar domínio: %w", err)
	}

	return result, nil
}

// UpdateDNS atualiza um domínio existente
func (s *DNSService) UpdateDNS(id int, req models.DNSUpdateRequest) (*models.DNSUpdateResponse, error) {
	// Na API da GoCache, a atualização de DNS é feita pelo ID do registro
	endpoint := fmt.Sprintf("/dns/%d", id)
	result := &models.DNSUpdateResponse{}

	_, err := s.Client.Put(endpoint, req, result)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar domínio: %w", err)
	}

	return result, nil
}

// DeleteDNS exclui um domínio pelo ID
func (s *DNSService) DeleteDNS(id int) (*models.DNSDeleteResponse, error) {
	// Na API da GoCache, a exclusão de DNS é feita pelo ID do registro
	endpoint := fmt.Sprintf("/dns/%d", id)
	result := &models.DNSDeleteResponse{}

	_, err := s.Client.Delete(endpoint, result)
	if err != nil {
		return nil, fmt.Errorf("erro ao excluir domínio: %w", err)
	}

	return result, nil
}
