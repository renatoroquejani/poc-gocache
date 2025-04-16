package models

// SmartRuleRewriteMatch representa as condições para ativar uma regra de redirecionamento
type SmartRuleRewriteMatch struct {
	RequestURI     string   `json:"request_uri,omitempty" form:"match[request_uri]"`
	Request        string   `json:"request,omitempty" form:"match[request]"` // Mantido para compatibilidade
	RequestMethods []string `json:"request_method,omitempty" form:"match[request_method]"`
	DeviceTypes    []string `json:"device_type,omitempty" form:"match[device_type]"`
	Host           string   `json:"host,omitempty" form:"match[host]"`
}

// SmartRuleRewriteAction representa a ação a ser executada quando a regra de redirecionamento é ativada
type SmartRuleRewriteAction struct {
	RedirectType string `json:"redirect_type,omitempty" form:"action[redirect_type]"`
	RedirectTo   string `json:"redirect_to,omitempty" form:"action[redirect_to]"`
	RewriteURI   string `json:"rewrite_uri,omitempty" form:"action[rewrite_uri]"`
	RewriteHost  string `json:"rewrite_host,omitempty" form:"action[rewrite_host]"`
	Destination  string `json:"destination,omitempty" form:"action[destination]"`
	CrossOrigin  string `json:"cross_origin,omitempty" form:"action[cross_origin]"`
}

// SmartRuleRewriteMetadata representa metadados adicionais da regra de redirecionamento
type SmartRuleRewriteMetadata struct {
	Status    string `json:"status,omitempty"`
	UpdatedOn string `json:"updated_on,omitempty"`
}

// SmartRuleRewrite representa uma regra de redirecionamento
type SmartRuleRewrite struct {
	Match    SmartRuleRewriteMatch    `json:"match"`
	Action   SmartRuleRewriteAction   `json:"action"`
	ID       string                   `json:"id,omitempty"`
	Metadata SmartRuleRewriteMetadata `json:"metadata,omitempty"`
}

// SmartRuleRewriteCreateRequest representa a requisição para criar uma nova regra de redirecionamento
type SmartRuleRewriteCreateRequest struct {
	Match  SmartRuleRewriteMatch  `json:"match"`
	Action SmartRuleRewriteAction `json:"action"`
	Domain string                 `json:"-"` // Campo para armazenar o domínio, não será serializado para JSON
}

// SmartRuleRewriteCreateResponse representa a resposta da API para criação de regra de redirecionamento
type SmartRuleRewriteCreateResponse struct {
	Response struct {
		ID string `json:"id"`
	} `json:"response"`
}

// SmartRuleRewriteListResponse representa a resposta da API para listagem de regras de redirecionamento
type SmartRuleRewriteListResponse struct {
	Response struct {
		Rules []SmartRuleRewrite `json:"rules"`
	} `json:"response"`
}

// SmartRuleRewriteDeleteResponse representa a resposta da API para exclusão de regra de redirecionamento
type SmartRuleRewriteDeleteResponse struct {
	Response struct {
		Msg string `json:"msg"`
	} `json:"response"`
}

// SmartRuleRewriteUpdateResponse representa a resposta da API para atualização de regra de redirecionamento
type SmartRuleRewriteUpdateResponse struct {
	Response struct {
		Msg string `json:"msg"`
	} `json:"response"`
}

// SmartRuleSimplifiedRequest representa uma requisição simplificada para criar uma regra padrão
type SmartRuleSimplifiedRequest struct {
	Domain       string `json:"domain" binding:"required"`        // Subdomínio (campo unificado com nome consistente)
	ParentDomain string `json:"parent_domain" binding:"required"` // Domínio principal já existente na GoCache (ex: sites.kodestech.com.br)
	BucketURL    string `json:"bucket_url" binding:"required"`    // URL do bucket (ex: onm-landing-pages.s3-website-us-east-1.amazonaws.com)
	AccountID    string `json:"account_id" binding:"required"`    // ID da conta (ex: cliente-1)
}

// SmartRuleSimplifiedFormResponse representa a resposta para o formulário de criação de regra simplificada
type SmartRuleSimplifiedFormResponse struct {
	Domains []DomainOption `json:"domains"` // Lista de domínios disponíveis
}

// DomainOption representa uma opção de domínio para seleção no formulário
type DomainOption struct {
	Name        string `json:"name"`         // Nome do domínio (exod.com.br)
	DisplayName string `json:"display_name"` // Nome para exibição
}
