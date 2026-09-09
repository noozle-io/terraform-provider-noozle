package client

import (
	"context"
	"fmt"
)

type typedOutletCreateRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Fields      map[string]any `json:"fields"`
	Enabled     bool           `json:"enabled"`
}

type typedOutletUpdateRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Fields      map[string]any `json:"fields"`
	Enabled     bool           `json:"enabled"`
}

type typedOutletResponse struct {
	ID                int64            `json:"id"`
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	Enabled           bool             `json:"enabled"`
	Config            map[string]any   `json:"config"`
	StreamURL         string           `json:"stream_url"`
	ApplyStatus       string           `json:"apply_status"`
	LastApplyError    string           `json:"last_apply_error"`
	ConfigGeneration  int64            `json:"config_generation"`
	AppliedGeneration int64            `json:"applied_generation"`
	Operation         *OutletOperation `json:"operation"`
}

type outletFieldMutationRequest struct {
	Value any `json:"value"`
}

type typedOutletDeleteResponse struct {
	Operation OutletOperation `json:"operation"`
}

type payloadFileRef struct {
	File string `json:"$file"`
}

func (c *Client) patchOutletField(ctx context.Context, path string, value any) error {
	return c.doJSON(ctx, httpMethodPatch, path, outletFieldMutationRequest{Value: value}, nil, statusOK)
}

func (c *Client) deleteOutletField(ctx context.Context, path string) error {
	return c.doJSON(ctx, httpMethodDelete, path, nil, nil, statusOK)
}

func (c *Client) deleteTypedOutlet(ctx context.Context, path string) (OutletOperation, error) {
	var response typedOutletDeleteResponse
	if err := c.doJSON(ctx, httpMethodDelete, path, nil, &response, statusOK, statusAccepted, statusNoContent); err != nil {
		return OutletOperation{}, err
	}

	return response.Operation, nil
}

func stringField(fields map[string]any, key string) (string, error) {
	if fields == nil {
		return "", nil
	}

	value, ok := fields[key]
	if !ok || value == nil {
		return "", nil
	}

	stringValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("expected %s to be a string, got %T", key, value)
	}

	return stringValue, nil
}

func requiredStringField(fields map[string]any, key string) (string, error) {
	value, err := stringField(fields, key)
	if err != nil {
		return "", err
	}
	if value == "" {
		return "", fmt.Errorf("missing required field %s in outlet config", key)
	}

	return value, nil
}

func requiredInt64Field(fields map[string]any, key string) (int64, error) {
	value, ok, err := int64Field(fields, key)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, fmt.Errorf("missing required field %s in outlet config", key)
	}

	return value, nil
}

func int64Field(fields map[string]any, key string) (int64, bool, error) {
	if fields == nil {
		return 0, false, nil
	}

	value, ok := fields[key]
	if !ok || value == nil {
		return 0, false, nil
	}

	switch typed := value.(type) {
	case float64:
		return int64(typed), true, nil
	case int:
		return int64(typed), true, nil
	case int64:
		return typed, true, nil
	default:
		return 0, false, fmt.Errorf("expected %s to be a number, got %T", key, value)
	}
}

func boolField(fields map[string]any, key string) (bool, bool, error) {
	if fields == nil {
		return false, false, nil
	}

	value, ok := fields[key]
	if !ok || value == nil {
		return false, false, nil
	}

	boolValue, ok := value.(bool)
	if !ok {
		return false, false, fmt.Errorf("expected %s to be a boolean, got %T", key, value)
	}

	return boolValue, true, nil
}

func float64Field(fields map[string]any, key string) (float64, bool, error) {
	if fields == nil {
		return 0, false, nil
	}

	value, ok := fields[key]
	if !ok || value == nil {
		return 0, false, nil
	}

	switch typed := value.(type) {
	case float64:
		return typed, true, nil
	case int:
		return float64(typed), true, nil
	case int64:
		return float64(typed), true, nil
	default:
		return 0, false, fmt.Errorf("expected %s to be a number, got %T", key, value)
	}
}

func stringListField(fields map[string]any, key string) ([]string, error) {
	if fields == nil {
		return nil, nil
	}

	value, ok := fields[key]
	if !ok || value == nil {
		return nil, nil
	}

	switch typed := value.(type) {
	case []string:
		return typed, nil
	case []any:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			stringValue, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("expected %s to contain strings, got %T", key, item)
			}
			values = append(values, stringValue)
		}
		return values, nil
	default:
		return nil, fmt.Errorf("expected %s to be a string list, got %T", key, value)
	}
}

func stringMapField(fields map[string]any, key string) (map[string]string, error) {
	if fields == nil {
		return nil, nil
	}

	value, ok := fields[key]
	if !ok || value == nil {
		return nil, nil
	}

	switch typed := value.(type) {
	case map[string]string:
		return typed, nil
	case map[string]any:
		values := make(map[string]string, len(typed))
		for mapKey, mapValue := range typed {
			stringValue, ok := mapValue.(string)
			if !ok {
				return nil, fmt.Errorf("expected %s to contain string values, got %T for key %s", key, mapValue, mapKey)
			}
			values[mapKey] = stringValue
		}
		return values, nil
	default:
		return nil, fmt.Errorf("expected %s to be an object, got %T", key, value)
	}
}

func stringFieldPointer(fields map[string]any, key string) (*string, error) {
	value, err := stringField(fields, key)
	if err != nil {
		return nil, err
	}
	if value == "" {
		return nil, nil
	}

	return &value, nil
}

// notifyPolicyPointerPreservingEmpty distinguishes an absent notify policy
// from an explicitly configured empty string. Terraform needs that distinction
// for optional fields whose empty value is meaningful configuration.
func notifyPolicyPointerPreservingEmpty(fields map[string]any) (*string, error) {
	if fields == nil {
		return nil, nil
	}

	value, ok := fields["notify_policy"]
	if !ok || value == nil {
		return nil, nil
	}

	stringValue, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("expected notify_policy to be a string, got %T", value)
	}

	return &stringValue, nil
}

func int64FieldPointer(fields map[string]any, key string) (*int64, error) {
	value, ok, err := int64Field(fields, key)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &value, nil
}

func boolFieldPointer(fields map[string]any, key string) (*bool, error) {
	value, ok, err := boolField(fields, key)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &value, nil
}

func float64FieldPointer(fields map[string]any, key string) (*float64, error) {
	value, ok, err := float64Field(fields, key)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &value, nil
}

func setOptionalStringField(fields map[string]any, key string, value *string) {
	if value != nil {
		fields[key] = *value
	}
}

func setOptionalInt64Field(fields map[string]any, key string, value *int64) {
	if value != nil {
		fields[key] = *value
	}
}

func setOptionalBoolField(fields map[string]any, key string, value *bool) {
	if value != nil {
		fields[key] = *value
	}
}

func setOptionalFloat64Field(fields map[string]any, key string, value *float64) {
	if value != nil {
		fields[key] = *value
	}
}

func setOptionalStringListField(fields map[string]any, key string, value []string) {
	if value != nil {
		fields[key] = value
	}
}

func setOptionalStringMapField(fields map[string]any, key string, value map[string]string) {
	if value != nil {
		fields[key] = value
	}
}
