# Guia de Uso - API de Integração Gocache

## Visão Geral

Esta API permite gerenciar domínios, smart rules e cache na CDN Gocache. A documentação interativa está disponível via Swagger.

## Executando a API

1. Certifique-se de ter Go 1.20+ instalado
2. Configure o arquivo `.env` com sua chave de API:
   ```
   GOCACHE_API_KEY=sua_chave_api_aqui
   GOCACHE_API_URL=https://api.gocache.com.br/v1
   PORT=8081
   ```
3. Execute o binário ou use o comando: `go run cmd/api/main.go`
4. Acesse a documentação Swagger: `http://localhost:8081/swagger/index.html`
5. Para o serviço de proxy: `go run cmd/proxy/main.go`

## Recursos Disponíveis

### DNS (Domínios)

* **Listar Registros DNS**
  - Endpoint: `GET /api/v1/dns?domain=example.com`
  - Descrição: Retorna todos os registros DNS cadastrados para um domínio específico

* **Obter Registro DNS**
  - Endpoint: `GET /api/v1/dns/{id}`
  - Descrição: Retorna detalhes de um registro DNS específico

* **Criar Registro DNS**
  - Endpoint: `POST /api/v1/dns/{domain}`
  - Descrição: Cria um novo registro DNS para um domínio
  - Corpo da requisição:
    ```json
    {
      "name": "www.example.com",
      "type": "CNAME",
      "content": "example.origin.com",
      "ttl": 3600,
      "cloud": true
    }
    ```

* **Atualizar Registro DNS**
  - Endpoint: `PUT /api/v1/dns/{id}`
  - Descrição: Atualiza um registro DNS existente
  - Corpo da requisição: Similar ao de criação

* **Remover Registro DNS**
  - Endpoint: `DELETE /api/v1/dns/{id}`
  - Descrição: Remove um registro DNS específico

### Smart Rules

* **Listar Smart Rules de Reescrita**
  - Endpoint: `GET /api/v1/rules/settings/{domain}`
  - Descrição: Lista todas as regras de reescrita de um domínio

* **Criar Smart Rule de Reescrita**
  - Endpoint: `POST /api/v1/rules/settings/{domain}`
  - Descrição: Cria uma nova regra de reescrita
  - Corpo da requisição:
    ```json
    {
      "name": "Minha Regra",
      "description": "Descrição da regra",
      "conditions": [
        {
          "type": "path",
          "operator": "matches",
          "value": "/pagina-*"
        }
      ],
      "actions": [
        {
          "type": "rewrite",
          "action_value": "/nova-pagina"
        }
      ]
    }
    ```

* **Atualizar Smart Rule de Reescrita**
  - Endpoint: `PUT /api/v1/rules/settings/{domain}/{id}`
  - Descrição: Atualiza uma regra de reescrita existente
  - Corpo da requisição: Similar ao de criação

* **Remover Smart Rule de Reescrita**
  - Endpoint: `DELETE /api/v1/rules/settings/{domain}/{id}`
  - Descrição: Remove uma regra de reescrita específica

* **Criar Regra Simplificada**
  - Endpoint: `POST /api/v1/rules/{domain}/simplified`
  - Descrição: Cria uma regra simplificada para reescrita ou redirecionamento
  - Corpo da requisição:
    ```json
    {
      "source_path": "/caminho-antigo",
      "target_path": "/caminho-novo",
      "rule_type": "rewrite",
      "description": "Redireciona caminho antigo para o novo"
    }
    ```

* **Obter Formulário de Regra Simplificada**
  - Endpoint: `GET /api/v1/rules/simplified/form`
  - Descrição: Obtém o formulário para criação de regra simplificada

### Cache

* **Expirar Cache por URLs**
  - Endpoint: `POST /api/v1/cache/purge`
  - Descrição: Expira o cache de URLs específicas
  - Corpo da requisição:
    ```json
    {
      "domain_id": 1234,
      "urls": [
        "https://example.com/pagina1.html",
        "https://example.com/pagina2.html"
      ]
    }
    ```

* **Expirar Cache por Prefixo**
  - Endpoint: `POST /api/v1/cache/purge-by-prefix`
  - Descrição: Expira o cache de todas as URLs com um prefixo específico
  - Corpo da requisição:
    ```json
    {
      "domain_id": 1234,
      "prefix": "/images/"
    }
    ```

* **Status do Cache**
  - Endpoint: `GET /api/v1/cache/status/{domainId}`
  - Descrição: Obtém o status atual do cache para um domínio

## Exemplos de Uso

### 1. Criar um registro DNS e configurar regra de reescrita

1. Primeiro, crie um registro DNS:
   ```
   POST /api/v1/dns/meusite.com
   {
     "name": "www",
     "type": "CNAME",
     "content": "meusite.origem.com",
     "ttl": 3600,
     "cloud": true
   }
   ```

2. Depois, configure uma regra de reescrita:
   ```
   POST /api/v1/rules/settings/meusite.com
   {
     "name": "Regra de Reescrita",
     "description": "Redireciona páginas antigas para novas",
     "conditions": [
       {
         "type": "path",
         "operator": "matches",
         "value": "/pagina-antiga/*"
       }
     ],
     "actions": [
       {
         "type": "rewrite",
         "action_value": "/pagina-nova/$1"
       }
     ]
   }
   ```

### 2. Expirar cache de URLs

```
POST /api/v1/cache/purge
{
  "domain_id": 1234,
  "urls": [
    "https://meusite.com/pagina-atualizada.html"
  ]
}
```

## Documentação de Referência

Para mais detalhes sobre a API do Gocache, consulte a documentação oficial:

- [API de Domínios](https://docs.gocache.com.br/api/api_dominios/)
- [API de Smart Rules](https://docs.gocache.com.br/api/api_smart_rules/)
- [API de Cache](https://docs.gocache.com.br/api/api_cache/)