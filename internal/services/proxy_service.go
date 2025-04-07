package services

import (
	"fmt"
	"log"
	"sync"

	"github.com/renatoroquejani/poc-gocache/internal/models"
)

// ProxyService gerencia os mapeamentos de domu00ednios para destinos
type ProxyService struct {
	mappings []models.DomainMapping
	mutex    sync.RWMutex
}

// NewProxyService cria uma nova instu00e2ncia do serviu00e7o de proxy
func NewProxyService() *ProxyService {
	// Inicializa com alguns mapeamentos pru00e9-definidos
	return &ProxyService{
		mappings: []models.DomainMapping{
			{
				Domain:      "elizio.sites.kodestech.com.br",
				Destination: "https://onm-funnel-builder-stg.s3.us-east-2.amazonaws.com/account_pages/bolo-brigadeiro/index.html",
			},
			// Adicione mais mapeamentos conforme necessidade
		},
	}
}

// AddMapping adiciona um novo mapeamento de domu00ednio
func (s *ProxyService) AddMapping(mapping models.DomainMapping) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Verifica se o domu00ednio ju00e1 existe
	for i, m := range s.mappings {
		if m.Domain == mapping.Domain {
			// Atualiza o mapeamento existente
			s.mappings[i] = mapping
			log.Printf("Mapeamento atualizado para o domu00ednio %s: %s", mapping.Domain, mapping.Destination)
			return nil
		}
	}

	// Adiciona novo mapeamento
	s.mappings = append(s.mappings, mapping)
	log.Printf("Novo mapeamento adicionado para o domu00ednio %s: %s", mapping.Domain, mapping.Destination)
	return nil
}

// GetMapping retorna o mapeamento para um domu00ednio especu00edfico
func (s *ProxyService) GetMapping(domain string) (models.DomainMapping, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, m := range s.mappings {
		if m.Domain == domain {
			return m, nil
		}
	}

	return models.DomainMapping{}, fmt.Errorf("mapeamento nu00e3o encontrado para o domu00ednio: %s", domain)
}

// GetAllMappings retorna todos os mapeamentos
func (s *ProxyService) GetAllMappings() []models.DomainMapping {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Cria uma cu00f3pia para evitar problemas de concorru00eancia
	mappingsCopy := make([]models.DomainMapping, len(s.mappings))
	copy(mappingsCopy, s.mappings)

	return mappingsCopy
}

// DeleteMapping remove um mapeamento de domu00ednio
func (s *ProxyService) DeleteMapping(domain string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, m := range s.mappings {
		if m.Domain == domain {
			// Remove o mapeamento
			s.mappings = append(s.mappings[:i], s.mappings[i+1:]...)
			log.Printf("Mapeamento removido para o domu00ednio %s", domain)
			return nil
		}
	}

	return fmt.Errorf("mapeamento nu00e3o encontrado para o domu00ednio: %s", domain)
}
