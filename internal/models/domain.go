package models

// DomainCreateRequest represents a request to create a new domain in GoCache
type DomainCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Origin      string `json:"origin" binding:"required"`
	Description string `json:"description,omitempty"`
	Enabled     bool   `json:"enabled"`
}

// DomainInfo representa informações básicas de um domínio
type DomainInfo struct {
	Name string `json:"name"`
}

// DomainListResponse representa a resposta da API para listagem de domínios
type DomainListResponse struct {
	StatusCode int `json:"status_code"`
	Response struct {
		Domains []string `json:"domains"`
		Size   int      `json:"size"`
		AutoDiscovery map[string]interface{} `json:"auto_discovery"`
	} `json:"response"`
}
