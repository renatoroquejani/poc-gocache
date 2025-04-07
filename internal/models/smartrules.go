package models

// SmartRuleListResponse representa a resposta da API para listagem de smart rules
type SmartRuleListResponse struct {
	Status bool        `json:"status"`
	Data   []SmartRule `json:"data"`
}

// SmartRule representa uma smart rule na Gocache
type SmartRule struct {
	ID          int    `json:"id"`
	DomainID    int    `json:"domain_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Origin      string `json:"origin"`
	Path        string `json:"path"`
	Status      bool   `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// SmartRuleAction represents an action in a smart rule
type SmartRuleAction struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// SmartRuleCreateRequest representa a requisição para criar uma nova smart rule
type SmartRuleCreateRequest struct {
	DomainID     int               `json:"domain_id" binding:"required"`
	Name         string            `json:"name" binding:"required"`
	Description  string            `json:"description"`
	Origin       string            `json:"origin" binding:"required"`
	Path         string            `json:"path" binding:"required"`
	Actions      []SmartRuleAction `json:"actions"`
	CustomDomain string            `json:"custom_domain"`
}

// S3SmartRuleRequest é um wrapper para criar smart rules apontando para S3
type S3SmartRuleRequest struct {
	DomainID     int    `json:"domain_id" binding:"required"`
	S3Bucket     string `json:"s3_bucket" binding:"required"`
	UserFolder   string `json:"user_folder" binding:"required"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	CustomDomain string `json:"custom_domain"`
}

// SmartRuleUpdateRequest representa a requisição para atualizar uma smart rule
type SmartRuleUpdateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Origin      string `json:"origin" binding:"required"`
	Path        string `json:"path" binding:"required"`
}

// SmartRuleResponse representa a resposta da API para operações com smart rules
type SmartRuleResponse struct {
	Status  bool      `json:"status"`
	Data    SmartRule `json:"data"`
	Message string    `json:"message"`
}

// SmartRuleDeleteResponse representa a resposta da API para exclusão de smart rule
type SmartRuleDeleteResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
