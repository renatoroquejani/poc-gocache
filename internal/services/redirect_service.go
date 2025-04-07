package services

import (
	"fmt"
	"log"

	"github.com/renatoroquejani/poc-gocache/internal/models"
	"github.com/renatoroquejani/poc-gocache/pkg/gocache"
)

// RedirectService gerencia as operações relacionadas a regras de redirecionamento
type RedirectService struct {
	client *gocache.Client
}

// NewRedirectService cria uma nova instância do serviço de redirecionamento
func NewRedirectService(client *gocache.Client) *RedirectService {
	return &RedirectService{
		client: client,
	}
}

// CreateRedirect cria uma nova regra de redirecionamento
func (s *RedirectService) CreateRedirect(request *models.RedirectCreateRequest) (*models.RedirectCreateResponse, error) {
	log.Printf("Criando regra de redirecionamento para o domínio %s: %s -> %s",
		request.Domain, request.Source, request.Destination)

	endpoint := fmt.Sprintf("/redirects/%s", request.Domain)
	response := &models.RedirectCreateResponse{}

	_, err := s.client.Post(endpoint, request, response)
	if err != nil {
		log.Printf("Erro ao criar regra de redirecionamento: %v", err)
		return nil, fmt.Errorf("erro ao criar regra de redirecionamento: %w", err)
	}

	return response, nil
}

// ListRedirects lista todas as regras de redirecionamento para um domínio
func (s *RedirectService) ListRedirects(domain string) (*models.RedirectListResponse, error) {
	log.Printf("Listando regras de redirecionamento para o domínio %s", domain)

	endpoint := fmt.Sprintf("/redirects/%s", domain)
	response := &models.RedirectListResponse{}

	_, err := s.client.Get(endpoint, response)
	if err != nil {
		log.Printf("Erro ao listar regras de redirecionamento: %v", err)
		return nil, fmt.Errorf("erro ao listar regras de redirecionamento: %w", err)
	}

	return response, nil
}

// DeleteRedirect exclui uma regra de redirecionamento
func (s *RedirectService) DeleteRedirect(domain string, id int) (*models.RedirectDeleteResponse, error) {
	log.Printf("Excluindo regra de redirecionamento %d do domínio %s", id, domain)

	endpoint := fmt.Sprintf("/redirects/%s/%d", domain, id)
	response := &models.RedirectDeleteResponse{}

	_, err := s.client.DeleteSimple(endpoint, response)
	if err != nil {
		log.Printf("Erro ao excluir regra de redirecionamento: %v", err)
		return nil, fmt.Errorf("erro ao excluir regra de redirecionamento: %w", err)
	}

	return response, nil
}
