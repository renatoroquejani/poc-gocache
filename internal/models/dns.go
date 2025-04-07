package models

// DNSListResponse representa a resposta da API para listagem de domínios
type DNSListResponse struct {
	StatusCode int `json:"status_code"`
	Response   struct {
		Records []struct {
			Name     string `json:"name"`
			Content  string `json:"content"`
			Type     string `json:"type"`
			TTL      string `json:"ttl"`
			Cloud    string `json:"cloud"`
			RecordID string `json:"record_id"`
		} `json:"records"`
	} `json:"response"`
}

// DNSEntry representa um registro de domínio na Gocache
type DNSEntry struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Origin      string `json:"origin"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Status      bool   `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// DNSCreateRequest representa a requisição para criar um novo domínio
type DNSCreateRequest struct {
	Name    string `json:"name" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Content string `json:"content" binding:"required"`
	TTL     int    `json:"ttl" binding:"required"`
	Cloud   int    `json:"cloud" binding:"required"`
	Domain  string `json:"-"` // Campo para armazenar o domínio, não será serializado para JSON
}

// DNSUpdateRequest representa a requisição para atualizar um domínio existente
type DNSUpdateRequest struct {
	Name    string `json:"name" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Content string `json:"content" binding:"required"`
	TTL     int    `json:"ttl" binding:"required"`
	Cloud   int    `json:"cloud" binding:"required"`
}

// DNSCreateResponse representa a resposta da API para criação de domínio
type DNSCreateResponse struct {
	StatusCode int `json:"status_code"`
	Response   struct {
		Records []struct {
			Name     string      `json:"name"`
			Content  string      `json:"content"`
			Type     string      `json:"type"`
			TTL      string      `json:"ttl"`
			Cloud    string      `json:"cloud"`
			RecordID interface{} `json:"record_id"` // Pode ser string ou número
		} `json:"records"`
	} `json:"response"`
}

// DNSUpdateResponse representa a resposta da API para atualização de domínio
type DNSUpdateResponse struct {
	StatusCode int `json:"status_code"`
	Response   struct {
		Records []struct {
			Name     string      `json:"name"`
			Content  string      `json:"content"`
			Type     string      `json:"type"`
			TTL      string      `json:"ttl"`
			Cloud    string      `json:"cloud"`
			RecordID interface{} `json:"record_id"` // Pode ser string ou número
		} `json:"records"`
	} `json:"response"`
}

// DNSDeleteResponse representa a resposta da API para exclusão de domínio
type DNSDeleteResponse struct {
	StatusCode int    `json:"status_code"`
	Response   string `json:"response"`
}
