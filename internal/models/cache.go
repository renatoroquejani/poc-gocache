package models

// CachePurgeRequest representa a requisição para expirar cache de URLs
type CachePurgeRequest struct {
	DomainID int      `json:"domain_id" binding:"required"`
	URLs     []string `json:"urls" binding:"required"`
}

// CachePurgeByPrefixRequest representa a requisição para expirar cache por prefixo
type CachePurgeByPrefixRequest struct {
	DomainID int    `json:"domain_id" binding:"required"`
	Prefix   string `json:"prefix" binding:"required"`
}

// CacheInvalidationResponse representa a resposta da API para invalidação de cache
type CacheInvalidationResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// CacheStatusResponse representa a resposta da API para status de cache
type CacheStatusResponse struct {
	Status bool        `json:"status"`
	Data   CacheStatus `json:"data"`
}

// CacheStatus contém informações sobre o status de cache
type CacheStatus struct {
	Total     int `json:"total"`
	Processed int `json:"processed"`
	Pending   int `json:"pending"`
}
