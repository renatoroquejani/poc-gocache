# POC Gocache Integration

API para integração com o serviço de CDN Gocache, permitindo gerenciar registros DNS, criar regras de reescrita, redirecionar domínios e expirar cache de rotas.

## Funcionalidades

- Gerenciamento de registros DNS
- Configuração de regras de reescrita e redirecionamento
- Interface simplificada para criação de regras
- Expiração de cache de rotas específicas
- Serviço de proxy para redirecionamento

## Requisitos

- Go 1.20+
- Chave de API da Gocache

## Configuração

1. Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:
```
GOCACHE_API_KEY=sua_chave_api_aqui
GOCACHE_API_URL=https://api.gocache.com.br/v1
PORT=8081
PROXY_PORT=8082
```

2. Execute a API principal:
```
go run cmd/api/main.go
```

3. Execute o serviço de proxy (opcional):
```
go run cmd/proxy/main.go
```

4. Acesse a documentação Swagger:
```
http://localhost:8081/swagger/index.html
```

## Estrutura do Projeto

- `cmd/api`: Ponto de entrada da API principal
- `cmd/proxy`: Servidor de proxy para redirecionamento de domínios
- `internal/handlers`: Handlers HTTP
- `internal/models`: Modelos de dados
- `internal/services`: Lógica de negócio
- `pkg/gocache`: Cliente para API da Gocache
- `docs`: Documentação do Swagger

## Para Mais Informações

Consulte o arquivo `GUIA_DE_USO.md` para instruções detalhadas sobre como usar a API.
