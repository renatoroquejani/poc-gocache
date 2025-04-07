package models

// RedirectRule representa uma regra de redirecionamento na GoCache
type RedirectRule struct {
	ID          int    `json:"id"`
	Domain      string `json:"domain" binding:"required"`
	Source      string `json:"source" binding:"required"`
	Destination string `json:"destination" binding:"required"`
	Type        int    `json:"type" binding:"required"` // 301 (permanente) ou 302 (temporário)
}

// RedirectCreateRequest representa a requisição para criar uma regra de redirecionamento
type RedirectCreateRequest struct {
	Domain      string `json:"domain" binding:"required"`
	Source      string `json:"source" binding:"required"`
	Destination string `json:"destination" binding:"required"`
	Type        int    `json:"type" binding:"required"` // 301 (permanente) ou 302 (temporário)
}

// RedirectCreateResponse representa a resposta da API para criação de regra de redirecionamento
type RedirectCreateResponse struct {
	StatusCode int    `json:"status_code"`
	Response   string `json:"response"`
}

// RedirectListResponse representa a resposta da API para listagem de regras de redirecionamento
type RedirectListResponse struct {
	StatusCode int            `json:"status_code"`
	Response   []RedirectRule `json:"response"`
}

// RedirectDeleteResponse representa a resposta da API para exclusão de regra de redirecionamento
type RedirectDeleteResponse struct {
	StatusCode int    `json:"status_code"`
	Response   string `json:"response"`
}
