package main

import (
	"fmt"
	"log"

	"github.com/renatoroquejani/poc-gocache/internal/models"
	"github.com/renatoroquejani/poc-gocache/internal/services"
	"github.com/renatoroquejani/poc-gocache/pkg/gocache"
)

func main() {
	// TODO: Replace with your actual GoCache API URL and API Key
	baseURL := "https://api.gocache.com.br"
	apiKey := "YOUR_GOCACHE_API_KEY"

	client, err := gocache.NewClient(baseURL, apiKey)
	if err != nil {
		log.Fatalf("Failed to create GoCache client: %v", err)
	}

	domainService := services.NewDomainService(client)

	req := models.DomainCreateRequest{
		Name:        "elizio.sites.exod.com.br",
		Origin:      "onm-landing-pages.s3.us-east-1.amazonaws.com",
		Description: "Test domain for cliente-2",
		Enabled:     true,
	}

	resp, err := domainService.CreateDomain(req)
	if err != nil {
		log.Fatalf("Failed to create domain: %v", err)
	}

	fmt.Printf("Domain created successfully: %+v\n", resp)
}
