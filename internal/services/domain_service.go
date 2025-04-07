package services

import (
	"fmt"

	"github.com/renatoroquejani/poc-gocache/internal/models"
	"github.com/renatoroquejani/poc-gocache/pkg/gocache"
)

// DomainService handles GoCache domain operations
type DomainService struct {
	client *gocache.Client
}

// NewDomainService creates a new DomainService
func NewDomainService(client *gocache.Client) *DomainService {
	return &DomainService{client: client}
}

// CreateDomain creates a new domain in GoCache
func (s *DomainService) CreateDomain(req models.DomainCreateRequest) (map[string]interface{}, error) {
	var result map[string]interface{}
	endpoint := fmt.Sprintf("/domain/%s", req.Name)

	formData := map[string]string{
		"cache_ttl":  "86400",
		"waf_status": "false",
		"cdn_mode":   "cname",
	}

	_, err := s.client.Post(endpoint, formData, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain: %w", err)
	}
	return result, nil
}

// DeleteDomain deletes a domain in GoCache
func (s *DomainService) DeleteDomain(domainID int) error {
	var result map[string]interface{}
	endpoint := fmt.Sprintf("/domains/%d", domainID)
	_, err := s.client.DeleteSimple(endpoint, &result)
	if err != nil {
		return fmt.Errorf("failed to delete domain: %w", err)
	}
	return nil
}

// ListDomains lista todos os domínios disponíveis na GoCache
func (s *DomainService) ListDomains() (*models.DomainListResponse, error) {
	var response models.DomainListResponse
	endpoint := "/domain"
	
	_, err := s.client.Get(endpoint, &response)
	if err != nil {
		return nil, fmt.Errorf("falha ao listar domínios: %w", err)
	}
	
	return &response, nil
}
