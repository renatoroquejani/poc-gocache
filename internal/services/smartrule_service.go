package services

import (
	"fmt"
	"strings"

	"github.com/renatoroquejani/poc-gocache/internal/models"
	"github.com/renatoroquejani/poc-gocache/pkg/gocache"
)

// SmartRuleService fornece métodos para interagir com a API de smart rules da Gocache
type SmartRuleService struct {
	client *gocache.Client
}

// NewSmartRuleService cria uma nova instância de SmartRuleService
func NewSmartRuleService(client *gocache.Client) *SmartRuleService {
	return &SmartRuleService{
		client: client,
	}
}

// ListSmartRules lista todas as smart rules para um domínio
func (s *SmartRuleService) ListSmartRules(domainID int) (*models.SmartRuleListResponse, error) {
	endpoint := fmt.Sprintf("/domains/%d/smart-rules", domainID)
	result := &models.SmartRuleListResponse{}

	_, err := s.client.Get(endpoint, result)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar smart rules: %w", err)
	}

	return result, nil
}

// GetSmartRule obtém detalhes de uma smart rule específica
func (s *SmartRuleService) GetSmartRule(domainID, ruleID int) (*models.SmartRuleResponse, error) {
	endpoint := fmt.Sprintf("/domains/%d/smart-rules/%d", domainID, ruleID)
	result := &models.SmartRuleResponse{}

	_, err := s.client.Get(endpoint, result)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter smart rule: %w", err)
	}

	return result, nil
}

// CreateSmartRule cria uma nova smart rule
func (s *SmartRuleService) CreateSmartRule(req models.SmartRuleCreateRequest) (*models.SmartRuleResponse, error) {
	endpoint := fmt.Sprintf("/domains/%d/smart-rules", req.DomainID)
	result := &models.SmartRuleResponse{}

	_, err := s.client.Post(endpoint, req, result)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar smart rule: %w", err)
	}

	return result, nil
}

// CreateS3SmartRule cria uma smart rule apontando para um bucket S3
func (s *SmartRuleService) CreateS3SmartRule(req models.S3SmartRuleRequest) (*models.SmartRuleResponse, error) {
	// Monta a URL de origem para o bucket S3
	s3Origin := fmt.Sprintf("%s/%s/", req.S3Bucket, req.UserFolder)

	// Se não houver nome, gera um nome baseado no bucket/pasta
	name := req.Name
	if name == "" {
		nameParts := strings.Split(req.UserFolder, "/")
		lastPart := nameParts[len(nameParts)-1]
		if lastPart == "" && len(nameParts) > 1 {
			lastPart = nameParts[len(nameParts)-2]
		}
		name = fmt.Sprintf("S3-%s", lastPart)
	}

	// Define as ações da smart rule conforme necessário
	actions := []models.SmartRuleAction{
		{
			Type:  "set_cors",
			Value: req.CustomDomain, // e.g., https://elizio.sites.kodestech.com.br
		},
		{
			Type:  "rewrite_uri",
			Value: fmt.Sprintf("/%s/$1", req.UserFolder), // e.g., /cliente-1/$1
		},
		{
			Type:  "set_host_header",
			Value: req.S3Bucket, // e.g., onm-landing-pages.s3.us-east-1.amazonaws.com
		},
		{
			Type:  "set_origin",
			Value: req.S3Bucket, // e.g., onm-landing-pages.s3.us-east-1.amazonaws.com
		},
	}

	// Cria a smart rule
	smartRule := models.SmartRuleCreateRequest{
		DomainID:    req.DomainID,
		Name:        name,
		Description: req.Description,
		Origin:      s3Origin,
		Path:        "/",
		Actions:     actions,
	}

	return s.CreateSmartRule(smartRule)
}

// UpdateSmartRule atualiza uma smart rule existente
func (s *SmartRuleService) UpdateSmartRule(domainID, ruleID int, req models.SmartRuleUpdateRequest) (*models.SmartRuleResponse, error) {
	endpoint := fmt.Sprintf("/domains/%d/smart-rules/%d", domainID, ruleID)
	result := &models.SmartRuleResponse{}

	_, err := s.client.Put(endpoint, req, result)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar smart rule: %w", err)
	}

	return result, nil
}

// DeleteSmartRule exclui uma smart rule
func (s *SmartRuleService) DeleteSmartRule(domainID, ruleID int) (*models.SmartRuleDeleteResponse, error) {
	endpoint := fmt.Sprintf("/domains/%d/smart-rules/%d", domainID, ruleID)
	result := &models.SmartRuleDeleteResponse{}

	_, err := s.client.DeleteSimple(endpoint, result)
	if err != nil {
		return nil, fmt.Errorf("erro ao excluir smart rule: %w", err)
	}

	return result, nil
}
