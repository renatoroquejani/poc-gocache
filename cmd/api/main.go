package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Importação dos docs gerados pelo Swagger
	_ "github.com/renatoroquejani/poc-gocache/docs"

	"github.com/renatoroquejani/poc-gocache/internal/handlers"
	"github.com/renatoroquejani/poc-gocache/internal/services"
	"github.com/renatoroquejani/poc-gocache/pkg/gocache"
)

// @title Gocache Integration API
// @version 1.0
// @description API para integração com o serviço de CDN Gocache
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /api/v1
func main() {
	// Carrega variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	// Obtém a porta da aplicação
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Obtém as configurações da API da Gocache
	apiKey := os.Getenv("GOCACHE_API_KEY")
	if apiKey == "" {
		log.Fatal("Variável de ambiente GOCACHE_API_KEY não definida")
	}

	apiURL := os.Getenv("GOCACHE_API_URL")
	if apiURL == "" {
		apiURL = "https://api.gocache.com.br/v1"
	}

	// Cria o cliente da API
	client, err := gocache.NewClient(apiURL, apiKey)
	if err != nil {
		log.Fatalf("Erro ao criar cliente da API: %v", err)
	}

	// Inicializa os serviços
	dnsService := services.NewDNSService(client)
	// smartRuleService removido - usando apenas smartRuleRewriteService
	domainService := services.NewDomainService(client)
	cacheService := services.NewCacheService(client)
	redirectService := services.NewRedirectService(client)
	smartRuleRewriteService := services.NewSmartRuleRewriteService(client)
	proxyService := services.NewProxyService()

	// Inicializa os handlers
	dnsHandler := handlers.NewDNSHandler(dnsService)
	// smartRuleHandler removido - usando apenas smartRuleRewriteHandler
	cacheHandler := handlers.NewCacheHandler(cacheService)
	redirectHandler := handlers.NewRedirectHandler(redirectService)
	smartRuleRewriteHandler := handlers.NewSmartRuleRewriteHandler(smartRuleRewriteService)
	proxyHandler := handlers.NewProxyHandler(proxyService)
	domainHandler := handlers.NewDomainHandler(domainService, nil)

	// Inicializa o router
	router := gin.Default()

	// Adiciona middleware de recuperação e logger
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Middleware para processar redirecionamentos de domínio
	router.Use(func(c *gin.Context) {
		// Verifica se é uma requisição para a API ou para o Swagger
		if strings.HasPrefix(c.Request.URL.Path, "/api/") ||
			strings.HasPrefix(c.Request.URL.Path, "/swagger/") {
			c.Next()
			return
		}

		// Tenta encontrar um mapeamento para o host
		host := c.Request.Host
		path := c.Request.URL.Path

		// Remove a porta do host, se presente
		if strings.Contains(host, ":") {
			host = strings.Split(host, ":")[0]
		}

		// Procura o mapeamento correspondente
		mapping, err := proxyService.GetMapping(host)
		if err != nil {
			// Se não encontrou mapeamento, continua o processamento normal
			c.Next()
			return
		}

		// Encontrou mapeamento, faz o redirecionamento
		destination := mapping.Destination

		// Se o path não for a raiz, adiciona ao destino
		if path != "/" {
			// Remove a barra inicial do path se o destino já terminar com barra
			if strings.HasSuffix(destination, "/") && strings.HasPrefix(path, "/") {
				path = path[1:]
			}
			destination = destination + path
		}

		log.Printf("Redirecionando %s%s para: %s", host, path, destination)
		c.Redirect(http.StatusMovedPermanently, destination)
		c.Abort()
	})

	// Configura as rotas da API
	apiGroup := router.Group("/api/v1")
	{
		// Registra as rotas dos handlers
		dnsHandler.RegisterRoutes(apiGroup)
		// smartRuleHandler removido - usando apenas smartRuleRewriteHandler
		cacheHandler.RegisterRoutes(apiGroup)
		redirectHandler.RegisterRoutes(router)           // Registra as rotas de redirecionamento
		smartRuleRewriteHandler.RegisterRoutes(apiGroup) // Registra as rotas de Smart Rules de redirecionamento no grupo de API
		proxyHandler.RegisterRoutes(router)              // Registra as rotas de proxy
		domainHandler.RegisterRoutes(apiGroup)
	}

	// Configura o Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Inicia o servidor
	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Servidor iniciado em http://localhost%s", serverAddr)
	log.Printf("Documentação Swagger disponível em http://localhost%s/swagger/index.html", serverAddr)
	log.Printf("Serviço de redirecionamento de domínios ativado")

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
