package models

// DomainMapping armazena o mapeamento entre domu00ednios e seus destinos
type DomainMapping struct {
	Domain      string `json:"domain" binding:"required"`
	Destination string `json:"destination" binding:"required"`
}

// DomainMappingResponse representa a resposta da API para operau00e7u00f5es de mapeamento de domu00ednio
type DomainMappingResponse struct {
	Success bool         `json:"success"`
	Mapping DomainMapping `json:"mapping,omitempty"`
	Error   string       `json:"error,omitempty"`
}

// DomainMappingsListResponse representa a resposta da API para listagem de mapeamentos
type DomainMappingsListResponse struct {
	Success   bool            `json:"success"`
	Mappings  []DomainMapping `json:"mappings"`
	Total     int             `json:"total"`
	Error     string          `json:"error,omitempty"`
}
