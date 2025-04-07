package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// DomainMapping armazena o mapeamento entre domínios e seus destinos
type DomainMapping struct {
	Domain      string
	Destination string
}

func main() {
	// Carrega variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	// Obtém a porta da aplicação
	port := os.Getenv("PROXY_PORT")
	if port == "" {
		port = "8082"
	}

	// Configuração inicial de mapeamentos
	mappings := []DomainMapping{
		{
			Domain:      "elizio.sites.kodestech.com.br",
			Destination: "https://onm-funnel-builder-stg.s3.us-east-2.amazonaws.com/account_pages/bolo-brigadeiro/index.html",
		},
		// Adicione mais mapeamentos conforme necessário
	}

	// Inicializa o router
	router := gin.Default()

	// Adiciona middleware de recuperação e logger
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Rota para adicionar novos mapeamentos via API
	router.POST("/api/mappings", func(c *gin.Context) {
		var newMapping DomainMapping
		if err := c.ShouldBindJSON(&newMapping); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Valida os dados
		if newMapping.Domain == "" || newMapping.Destination == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "domínio e destino são obrigatórios"})
			return
		}

		// Adiciona o novo mapeamento
		mappings = append(mappings, newMapping)
		c.JSON(http.StatusCreated, newMapping)
	})

	// Rota para listar mapeamentos
	router.GET("/api/mappings", func(c *gin.Context) {
		c.JSON(http.StatusOK, mappings)
	})

	// Rota para processar redirecionamentos
	router.GET("/*path", func(c *gin.Context) {
		host := c.Request.Host
		path := c.Param("path")

		log.Printf("Recebida requisição para host: %s, path: %s", host, path)

		// Remove a porta do host, se presente
		if strings.Contains(host, ":") {
			host = strings.Split(host, ":")[0]
		}

		// Procura o mapeamento correspondente
		for _, mapping := range mappings {
			if mapping.Domain == host {
				destination := mapping.Destination
				
				// Se o path não for a raiz, adiciona ao destino
				if path != "/" {
					// Remove a barra inicial do path se o destino já terminar com barra
					if strings.HasSuffix(destination, "/") && strings.HasPrefix(path, "/") {
						path = path[1:]
					}
					destination = destination + path
				}

				log.Printf("Redirecionando para: %s", destination)
				c.Redirect(http.StatusMovedPermanently, destination)
				return
			}
		}

		// Se não encontrou mapeamento, retorna erro
		c.JSON(http.StatusNotFound, gin.H{"error": "domínio não configurado"})
	})

	// Inicia o servidor
	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Servidor de proxy iniciado em http://localhost%s", serverAddr)
	
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
