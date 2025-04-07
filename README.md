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

## Cenários de Uso para o Projeto ONM

Para o projeto ONM, temos 2 cenários de configuração:

### 1. Domínio já existente

Se o domínio já estiver configurado na Gocache, apenas siga para os passos de criação da rota DNS e Smart Rule.

### 2. Domínio customizado (novo)

Caso o domínio seja novo, será necessário criá-lo antes de adicionar a rota DNS. Importante notar que o endpoint da GoCache para criação de domínio é singular (domain) e não plural (domains):

1. **Criar o domínio**:
   ```
   POST /api/v1/domain/{nome-do-dominio}
   ```

2. **Criar a rota DNS** apontando do subdomínio para a origem (S3 ou site):
   ```
   POST /api/v1/dns/{dominio.com.br}
   ```

3. **Criar a Smart Rule** que faz o rewrite da URL para o accountID correto no S3:
   ```
   POST /api/v1/rules/{dominio.com.br}/simplified
   ```

A rota Simplified já cria as 4 configurações padrão para o funcionamento completo da rota.

### Informações Importantes

#### Gerenciamento de IDs

Cada recurso criado na Gocache gera um ID único:

- Cada domínio criado gera um ID
- Cada rota DNS criada gera um record ID
- Cada regra (rule) criada gera um ID

É fundamental armazenar esses IDs em um banco de dados para facilitar a manutenção futura, atualizações ou exclusões dos recursos.

#### Configuração de DNS no Provedor

O usuário precisará configurar o DNS no seu provedor para apontar para a Gocache:

- A Gocache gera um domínio próprio para configuração no CNAME
- O formato desse domínio é: `subdominio.dominio.cdn.gocache.net`
- O cliente deve criar um registro CNAME no seu provedor de DNS apontando o subdomínio para este endereço da Gocache

## Para Mais Informações

Consulte o arquivo `GUIA_DE_USO.md` para instruções detalhadas sobre como usar a API.
