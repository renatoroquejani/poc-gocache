package services

import (
	"fmt"

	"github.com/renatoroquejani/poc-gocache/internal/models"
	"github.com/renatoroquejani/poc-gocache/pkg/gocache"
)

// CacheService fornece métodos para interagir com a API de cache da Gocache
type CacheService struct {
	client *gocache.Client
}

// NewCacheService cria uma nova instância de CacheService
func NewCacheService(client *gocache.Client) *CacheService {
	return &CacheService{
		client: client,
	}
}

// PurgeCache expira o cache para URLs específicas
func (s *CacheService) PurgeCache(req models.CachePurgeRequest) (*models.CacheInvalidationResponse, error) {
	endpoint := fmt.Sprintf("/domains/%d/cache/purge", req.DomainID)
	result := &models.CacheInvalidationResponse{}

	body := map[string][]string{
		"urls": req.URLs,
	}

	_, err := s.client.Post(endpoint, body, result)
	if err != nil {
		return nil, fmt.Errorf("erro ao expirar cache: %w", err)
	}

	return result, nil
}

// PurgeCacheByPrefix expira o cache para URLs que começam com um prefixo
func (s *CacheService) PurgeCacheByPrefix(req models.CachePurgeByPrefixRequest) (*models.CacheInvalidationResponse, error) {
	endpoint := fmt.Sprintf("/domains/%d/cache/purge-by-prefix", req.DomainID)
	result := &models.CacheInvalidationResponse{}

	body := map[string]string{
		"prefix": req.Prefix,
	}

	_, err := s.client.Post(endpoint, body, result)
	if err != nil {
		return nil, fmt.Errorf("erro ao expirar cache por prefixo: %w", err)
	}

	return result, nil
}

// GetCacheStatus obtém o status atual do cache para um domínio
func (s *CacheService) GetCacheStatus(domainID int) (*models.CacheStatusResponse, error) {
	endpoint := fmt.Sprintf("/domains/%d/cache/status", domainID)
	result := &models.CacheStatusResponse{}

	_, err := s.client.Get(endpoint, result)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter status do cache: %w", err)
	}

	return result, nil
}
