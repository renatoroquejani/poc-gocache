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

// PurgeAllCache expira todo o cache de um domínio
func (s *CacheService) PurgeAllCache(domain string) (*models.CacheInvalidationResponse, error) {
	// Na API GoCache, usa-se a rota /cache/{dominio}/all para expurgar todo o cache
	endpoint := fmt.Sprintf("/cache/%s/all", domain)
	result := &models.CacheInvalidationResponse{}

	// Para expurgar todo o cache, enviamos um DELETE sem body
	_, err := s.client.DeleteSimple(endpoint, result)
	if err != nil {
		return nil, fmt.Errorf("erro ao expirar todo o cache: %w", err)
	}

	return result, nil
}

// PurgeUrls expira o cache para URLs específicas, podendo incluir máscaras/wildcards
func (s *CacheService) PurgeUrls(req models.CachePurgeRequest) (*models.CacheInvalidationResponse, error) {
	// Na API GoCache, o domínio é parte da URL
	endpoint := fmt.Sprintf("/cache/%s", req.Domain)
	result := &models.CacheInvalidationResponse{}

	// Prepara os dados para a requisição
	body := map[string]string{
		"content-type": "*", // Por padrão, usamos wildcard para limpar todos os content-types
	}
	
	// Adiciona cada URL como um item separado no formato urls[0], urls[1], etc.
	// As URLs podem conter wildcards (ex: http://example.com/blog/*)
	for i, url := range req.URLs {
		body[fmt.Sprintf("urls[%d]", i)] = url
	}

	_, err := s.client.Delete(endpoint, body, result)
	if err != nil {
		return nil, fmt.Errorf("erro ao expirar cache para URLs: %w", err)
	}

	return result, nil
}


