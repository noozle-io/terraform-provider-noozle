package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	defaultRetryMax     = 2
	defaultRetryWaitMin = 100 * time.Millisecond
	defaultRetryWaitMax = 500 * time.Millisecond
)

// ErrMissingTenantContext indicates an operation needs a configured tenant.
var ErrMissingTenantContext = errors.New("tenant context is required; set provider tenant_id or NOOZLE_TENANT_ID")

type TenantContext struct {
	TenantID       string
	OrganizationID string
}

type Option func(*config)

type config struct {
	httpClient    *http.Client
	logger        Logger
	retryMax      int
	tenantContext TenantContext
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(cfg *config) {
		cfg.httpClient = httpClient
	}
}

func WithRetryMax(retryMax int) Option {
	return func(cfg *config) {
		cfg.retryMax = retryMax
	}
}

func WithLogger(logger Logger) Option {
	return func(cfg *config) {
		cfg.logger = logger
	}
}

// WithTenantContext configures the tenant ownership context sent with every request.
// Tenant-owned operations additionally require a non-empty tenant ID before they run.
func WithTenantContext(tenantContext TenantContext) Option {
	return func(cfg *config) {
		cfg.tenantContext = TenantContext{
			TenantID:       strings.TrimSpace(tenantContext.TenantID),
			OrganizationID: strings.TrimSpace(tenantContext.OrganizationID),
		}
	}
}

type Logger interface {
	Debug(context.Context, string, map[string]any)
}

type multipartFilePart struct {
	FieldName string
	FileName  string
	Content   []byte
}

type Client struct {
	baseURL       *url.URL
	apiKey        string
	tenantContext TenantContext
	client        *retryablehttp.Client
	logger        Logger
}

type APIError struct {
	StatusCode    int
	Message       string
	Type          string
	CorrelationID string
	Details       map[string]any
	Method        string
	Path          string
}

func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}

	message := e.Message
	if message == "" {
		message = http.StatusText(e.StatusCode)
	}

	if e.CorrelationID == "" {
		return fmt.Sprintf("%s %s: %s", e.Method, e.Path, message)
	}

	return fmt.Sprintf("%s %s: %s (correlation_id=%s)", e.Method, e.Path, message, e.CorrelationID)
}

func New(baseURL, apiKey string, opts ...Option) (*Client, error) {
	if strings.TrimSpace(baseURL) == "" {
		return nil, errors.New("base URL is required")
	}
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("api key is required")
	}

	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base URL: %w", err)
	}
	if parsedBaseURL.Scheme == "" || parsedBaseURL.Host == "" {
		return nil, fmt.Errorf("base URL must include scheme and host: %q", baseURL)
	}

	parsedBaseURL.Path = strings.TrimRight(parsedBaseURL.Path, "/")
	parsedBaseURL.RawPath = strings.TrimRight(parsedBaseURL.RawPath, "/")

	cfg := config{
		retryMax: defaultRetryMax,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = cfg.retryMax
	retryClient.RetryWaitMin = defaultRetryWaitMin
	retryClient.RetryWaitMax = defaultRetryWaitMax
	retryClient.Logger = nil
	retryClient.CheckRetry = checkRetry

	if cfg.httpClient != nil {
		retryClient.HTTPClient = cfg.httpClient
	}

	return &Client{
		baseURL:       parsedBaseURL,
		apiKey:        apiKey,
		tenantContext: cfg.tenantContext,
		client:        retryClient,
		logger:        cfg.logger,
	}, nil
}

func (c *Client) hasTenantContext() bool {
	return c.tenantContext.TenantID != ""
}

func (c *Client) requireTenantContext() error {
	if !c.hasTenantContext() {
		return ErrMissingTenantContext
	}

	return nil
}

func (c *Client) setHeaders(header http.Header) {
	header.Set("Accept", "application/json")
	header.Set("X-API-Key", c.apiKey)
}

func (c *Client) setTenantContextHeaders(header http.Header) {
	if c.tenantContext.TenantID != "" {
		header.Set("X-Noozle-Tenant", c.tenantContext.TenantID)
	}
	if c.tenantContext.OrganizationID != "" {
		header.Set("X-Noozle-Organization", c.tenantContext.OrganizationID)
	}
}

func IsBadRequest(err error) bool {
	return hasStatus(err, http.StatusBadRequest)
}

func IsUnauthorized(err error) bool {
	return hasStatus(err, http.StatusUnauthorized)
}

func IsNotFound(err error) bool {
	return hasStatus(err, http.StatusNotFound)
}

func IsServerError(err error) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}

	return apiErr.StatusCode >= http.StatusInternalServerError
}

func (c *Client) doJSON(ctx context.Context, method, path string, requestBody any, responseBody any, expectedStatus ...int) error {
	return c.doTenantJSON(ctx, method, path, requestBody, responseBody, expectedStatus...)
}

func (c *Client) doTenantJSON(ctx context.Context, method, path string, requestBody any, responseBody any, expectedStatus ...int) error {
	if err := c.requireTenantContext(); err != nil {
		return err
	}

	return c.doJSONWithTenantContext(ctx, method, path, requestBody, responseBody, expectedStatus...)
}

func (c *Client) doJSONWithTenantContext(ctx context.Context, method, path string, requestBody any, responseBody any, expectedStatus ...int) error {
	var bodyReader io.Reader
	var requestPayload []byte
	if requestBody != nil {
		payload, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		requestPayload = payload
		bodyReader = bytes.NewReader(payload)
	}

	request, err := retryablehttp.NewRequestWithContext(ctx, method, c.urlFor(path), bodyReader)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	c.setHeaders(request.Header)
	c.setTenantContextHeaders(request.Header)
	if requestBody != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	c.logPayload(ctx, "noozle client request", method, path, 0, requestPayload)

	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, path, err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusNoContent {
		c.logPayload(ctx, "noozle client response", method, path, response.StatusCode, nil)
		if expectedHTTPStatus(response.StatusCode, expectedStatus) {
			return nil
		}
		return decodeAPIErrorPayload(response.StatusCode, nil, method, path)
	}

	responsePayload, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	c.logPayload(ctx, "noozle client response", method, path, response.StatusCode, responsePayload)

	if !expectedHTTPStatus(response.StatusCode, expectedStatus) {
		return decodeAPIErrorPayload(response.StatusCode, responsePayload, method, path)
	}
	recordOutletMutationOperation(ctx, method, responsePayload)

	if responseBody == nil {
		return nil
	}

	if err := json.Unmarshal(responsePayload, responseBody); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}

	return nil
}

func (c *Client) doMultipartPayload(ctx context.Context, method, path string, payload any, files []multipartFilePart, responseBody any, expectedStatus ...int) error {
	if err := c.requireTenantContext(); err != nil {
		return err
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	if err := writer.WriteField("payload", string(payloadBytes)); err != nil {
		return fmt.Errorf("write payload field: %w", err)
	}

	for _, file := range files {
		part, err := writer.CreateFormFile(file.FieldName, filepath.Base(file.FileName))
		if err != nil {
			return fmt.Errorf("create file part %s: %w", file.FieldName, err)
		}

		if _, err := part.Write(file.Content); err != nil {
			return fmt.Errorf("write file part %s: %w", file.FieldName, err)
		}
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("close multipart writer: %w", err)
	}

	request, err := retryablehttp.NewRequestWithContext(ctx, method, c.urlFor(path), bytes.NewReader(body.Bytes()))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	c.setHeaders(request.Header)
	c.setTenantContextHeaders(request.Header)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	c.logPayload(ctx, "noozle client request", method, path, 0, body.Bytes())

	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, path, err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusNoContent {
		c.logPayload(ctx, "noozle client response", method, path, response.StatusCode, nil)
		if expectedHTTPStatus(response.StatusCode, expectedStatus) {
			return nil
		}
		return decodeAPIErrorPayload(response.StatusCode, nil, method, path)
	}

	responsePayload, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	c.logPayload(ctx, "noozle client response", method, path, response.StatusCode, responsePayload)

	if !expectedHTTPStatus(response.StatusCode, expectedStatus) {
		return decodeAPIErrorPayload(response.StatusCode, responsePayload, method, path)
	}
	recordOutletMutationOperation(ctx, method, responsePayload)

	if responseBody == nil {
		return nil
	}

	if err := json.Unmarshal(responsePayload, responseBody); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}

	return nil
}

func decodeBase64File(contentBase64 string) ([]byte, error) {
	content, err := base64.StdEncoding.DecodeString(contentBase64)
	if err != nil {
		return nil, fmt.Errorf("decode base64 content: %w", err)
	}

	return content, nil
}

func (c *Client) urlFor(path string) string {
	return c.baseURL.String() + path
}

func expectedHTTPStatus(statusCode int, expectedStatus []int) bool {
	for _, expected := range expectedStatus {
		if statusCode == expected {
			return true
		}
	}

	return false
}

func hasStatus(err error, statusCode int) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}

	return apiErr.StatusCode == statusCode
}

func checkRetry(_ context.Context, resp *http.Response, err error) (bool, error) {
	if hasRequestError(err) {
		return true, nil
	}
	if resp == nil || resp.Request == nil {
		return false, nil
	}
	if resp.Request.Method != http.MethodGet {
		return false, nil
	}
	if resp.StatusCode >= http.StatusInternalServerError {
		return true, nil
	}

	return false, nil
}

func hasRequestError(err error) bool {
	return err != nil
}

func decodeAPIErrorPayload(statusCode int, body []byte, method, path string) error {
	apiErr := &APIError{
		StatusCode: statusCode,
		Message:    http.StatusText(statusCode),
		Method:     method,
		Path:       path,
	}

	if len(body) == 0 {
		return apiErr
	}

	var payload errorResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return apiErr
	}

	if payload.Error != "" {
		apiErr.Message = payload.Error
	}
	apiErr.Type = payload.Type
	apiErr.CorrelationID = payload.CorrelationID
	apiErr.Details = payload.Details

	return apiErr
}

func (c *Client) logPayload(ctx context.Context, message, method, path string, statusCode int, payload []byte) {
	if c.logger == nil {
		return
	}

	fields := map[string]any{
		"method":        method,
		"path":          path,
		"payload_bytes": len(payload),
	}
	if statusCode != 0 {
		fields["status_code"] = statusCode
	}

	c.logger.Debug(ctx, message, fields)
}

type errorResponse struct {
	Error         string         `json:"error"`
	Type          string         `json:"type"`
	CorrelationID string         `json:"correlation_id"`
	Details       map[string]any `json:"details"`
}
