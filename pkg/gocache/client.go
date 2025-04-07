package gocache

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client representa o cliente para a API da Gocache
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *resty.Client
}

// NewClient cria uma nova instância do cliente da API da Gocache
func NewClient(baseURL, apiKey string) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("baseURL não pode ser vazia")
	}
	if apiKey == "" {
		return nil, errors.New("apiKey não pode ser vazia")
	}

	httpClient := resty.New()
	httpClient.SetTimeout(30 * time.Second)
	httpClient.SetRetryCount(3)
	httpClient.SetRetryWaitTime(5 * time.Second)
	httpClient.SetRetryMaxWaitTime(20 * time.Second)
	
	// Adiciona logger para debug
	httpClient.SetDebug(true)
	httpClient.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		log.Printf("Enviando requisição: %s %s", req.Method, req.URL)
		return nil
	})
	
	httpClient.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		log.Printf("Resposta: %d %s", resp.StatusCode(), resp.Status())
		log.Printf("Corpo da Resposta: %s", resp.String())
		return nil
	})

	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: httpClient,
	}, nil
}

// setAuthHeaders adiciona os headers de autenticação para as requisições
func (c *Client) setAuthHeaders(req *resty.Request) *resty.Request {
	// Formato conforme documentação da API Gocache
	return req.SetHeader("GoCache-Token", c.apiKey)
}

// Get realiza uma requisição GET para a API do Gocache
func (c *Client) Get(endpoint string, result interface{}) (*resty.Response, error) {
	return c.GetWithQueryParams(endpoint, nil, result)
}

// GetWithQueryParams realiza uma requisição GET com query parameters para a API do Gocache
func (c *Client) GetWithQueryParams(endpoint string, queryParams map[string]string, result interface{}) (*resty.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)
	req := c.httpClient.R().
		SetResult(result).
		EnableTrace()
	
	req = c.setAuthHeaders(req)
	
	// Adiciona query parameters se fornecidos
	if queryParams != nil {
		req.SetQueryParams(queryParams)
	}
	
	log.Printf("Enviando GET para: %s", url)
	log.Printf("Headers: %v", req.Header)
	if len(queryParams) > 0 {
		log.Printf("Query Params: %v", queryParams)
	}
	log.Printf("API Key: %s (primeiros 5 caracteres)", c.apiKey[:5])
	
	resp, err := req.Get(url)
	
	if err != nil {
		log.Printf("Erro na requisição GET: %v", err)
	} else {
		log.Printf("Resposta status code: %d", resp.StatusCode())
		log.Printf("Resposta corpo: %s", resp.String())
	}
	
	return resp, err
}

// Post realiza uma requisição POST para a API do Gocache
func (c *Client) Post(endpoint string, body, result interface{}) (*resty.Response, error) {
	return c.doRequest("POST", endpoint, body, result)
}

// doRequest realiza uma requisição genérica para a API do Gocache
func (c *Client) doRequest(method, endpoint string, body, result interface{}) (*resty.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)
	req := c.httpClient.R().
		SetResult(result).
		EnableTrace()
	
	// Adiciona os headers de autenticação
	c.setAuthHeaders(req)
	
	// Verifica se o body é um map[string]string e o usa diretamente como form data
	if formData, ok := body.(map[string]string); ok {
		// Se for um map[string]string, usa diretamente como form data
		req.SetFormData(formData)
	} else {
		// Caso contrário, tenta converter para map[string]string usando reflection
		formData := make(map[string]string)
		
		// Usa reflection para extrair os campos e valores do body
		val := reflect.ValueOf(body)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		
		if val.Kind() != reflect.Struct {
			return nil, fmt.Errorf("body deve ser uma struct, um ponteiro para struct ou um map[string]string")
		}
		
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			
			// Ignora campos privados
			if !field.IsExported() {
				continue
			}
			
			// Obtem o nome do campo para o formulário a partir da tag "form"
			formName := field.Tag.Get("form")
			if formName == "" {
				// Se não tiver tag "form", usa a tag "json"
				formName = field.Tag.Get("json")
				if formName == "" {
					// Se não tiver tag "json", usa o nome do campo
					formName = field.Name
				} else {
					// Remove opções da tag json (ex: omitempty)
					formName = strings.Split(formName, ",")[0]
				}
			}
			
			// Ignora campos marcados com "-"
			if formName == "-" {
				continue
			}
			
			fieldValue := val.Field(i)
			
			// Processa o valor do campo de acordo com seu tipo
			switch fieldValue.Kind() {
			case reflect.String:
				if fieldValue.String() != "" {
					formData[formName] = fieldValue.String()
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				formData[formName] = strconv.FormatInt(fieldValue.Int(), 10)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				formData[formName] = strconv.FormatUint(fieldValue.Uint(), 10)
			case reflect.Float32, reflect.Float64:
				formData[formName] = strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64)
			case reflect.Bool:
				formData[formName] = strconv.FormatBool(fieldValue.Bool())
			case reflect.Slice, reflect.Array:
				// Para slices e arrays, adiciona cada elemento como um item separado
				for j := 0; j < fieldValue.Len(); j++ {
					item := fieldValue.Index(j)
					switch item.Kind() {
					case reflect.String:
						formData[fmt.Sprintf("%s[%d]", formName, j)] = item.String()
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						formData[fmt.Sprintf("%s[%d]", formName, j)] = strconv.FormatInt(item.Int(), 10)
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						formData[fmt.Sprintf("%s[%d]", formName, j)] = strconv.FormatUint(item.Uint(), 10)
					case reflect.Float32, reflect.Float64:
						formData[fmt.Sprintf("%s[%d]", formName, j)] = strconv.FormatFloat(item.Float(), 'f', -1, 64)
					case reflect.Bool:
						formData[fmt.Sprintf("%s[%d]", formName, j)] = strconv.FormatBool(item.Bool())
					}
				}
			}
		}
		
		// Define os dados do formulário na requisição
		req.SetFormData(formData)
	}
	
	// Faz a requisição com o método especificado
	var resp *resty.Response
	var err error
	
	switch method {
	case "POST":
		resp, err = req.Post(url)
	case "DELETE":
		resp, err = req.Delete(url)
	case "PUT":
		resp, err = req.Put(url)
	case "PATCH":
		resp, err = req.Patch(url)
	default:
		return nil, fmt.Errorf("método HTTP não suportado: %s", method)
	}
	
	if err != nil {
		return nil, err
	}
	
	return resp, nil
}

// Put realiza uma requisição PUT para a API do Gocache
func (c *Client) Put(endpoint string, body, result interface{}) (*resty.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)
	req := c.httpClient.R().
		SetResult(result).
		EnableTrace()
	
	// Adiciona os headers de autenticação
	c.setAuthHeaders(req)
	
	// Verifica se o body é um map[string]string e o usa diretamente como form data
	if formData, ok := body.(map[string]string); ok {
		// Se for um map[string]string, usa diretamente como form data
		req.SetFormData(formData)
	} else {
		// Caso contrário, tenta converter para map[string]string usando reflection
		formData := make(map[string]string)
		
		// Usa reflection para extrair os campos e valores do body
		val := reflect.ValueOf(body)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		
		if val.Kind() != reflect.Struct {
			return nil, fmt.Errorf("body deve ser uma struct, um ponteiro para struct ou um map[string]string")
		}
		
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			
			// Ignora campos privados
			if !field.IsExported() {
				continue
			}
			
			// Obtem o nome do campo para o formulário a partir da tag "form"
			formName := field.Tag.Get("form")
			if formName == "" {
				// Se não tiver tag "form", usa a tag "json"
				formName = field.Tag.Get("json")
				if formName == "" {
					// Se não tiver tag "json", usa o nome do campo
					formName = field.Name
				} else {
					// Remove opções da tag json (ex: omitempty)
					formName = strings.Split(formName, ",")[0]
				}
			}
			
			// Ignora campos marcados com "-"
			if formName == "-" {
				continue
			}
			
			fieldValue := val.Field(i)
			
			// Processa o valor do campo de acordo com seu tipo
			switch fieldValue.Kind() {
			case reflect.String:
				if fieldValue.String() != "" {
					formData[formName] = fieldValue.String()
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				formData[formName] = strconv.FormatInt(fieldValue.Int(), 10)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				formData[formName] = strconv.FormatUint(fieldValue.Uint(), 10)
			case reflect.Float32, reflect.Float64:
				formData[formName] = strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64)
			case reflect.Bool:
				formData[formName] = strconv.FormatBool(fieldValue.Bool())
			case reflect.Slice, reflect.Array:
				// Para slices e arrays, adiciona cada elemento como um item separado
				for j := 0; j < fieldValue.Len(); j++ {
					item := fieldValue.Index(j)
					switch item.Kind() {
					case reflect.String:
						formData[fmt.Sprintf("%s[%d]", formName, j)] = item.String()
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						formData[fmt.Sprintf("%s[%d]", formName, j)] = strconv.FormatInt(item.Int(), 10)
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						formData[fmt.Sprintf("%s[%d]", formName, j)] = strconv.FormatUint(item.Uint(), 10)
					case reflect.Float32, reflect.Float64:
						formData[fmt.Sprintf("%s[%d]", formName, j)] = strconv.FormatFloat(item.Float(), 'f', -1, 64)
					case reflect.Bool:
						formData[fmt.Sprintf("%s[%d]", formName, j)] = strconv.FormatBool(item.Bool())
					}
				}
			}
		}
		
		// Define os dados do formulário na requisição
		req.SetFormData(formData)
	}
	
	// Faz a requisição PUT
	resp, err := req.Put(url)
	if err != nil {
		return nil, err
	}
	
	return resp, nil
}

// Delete realiza uma requisição DELETE para a API do Gocache
func (c *Client) Delete(endpoint string, body, result interface{}) (*resty.Response, error) {
	return c.doRequest("DELETE", endpoint, body, result)
}

// DeleteSimple realiza uma requisição DELETE simples sem body para a API do Gocache
func (c *Client) DeleteSimple(endpoint string, result interface{}) (*resty.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)
	req := c.httpClient.R().
		SetResult(result).
		EnableTrace()
	
	req = c.setAuthHeaders(req)
	resp, err := req.Delete(url)
	
	if err != nil {
		log.Printf("Erro na requisição DELETE: %v", err)
	}
	
	return resp, err
}
