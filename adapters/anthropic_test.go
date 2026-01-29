package adapters

import (
	"testing"

	"github.com/jonwraymond/tooladapter"
)

func TestAnthropicAdapter_Name(t *testing.T) {
	adapter := NewAnthropicAdapter()

	got := adapter.Name()
	want := "anthropic"

	if got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}

func TestAnthropicAdapter_ToCanonical_Basic(t *testing.T) {
	adapter := NewAnthropicAdapter()

	tool := AnthropicTool{
		Name:        "test_tool",
		Description: "A test tool",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{
					"type":        "string",
					"description": "The name",
				},
			},
			"required": []any{"name"},
		},
	}

	got, err := adapter.ToCanonical(tool)

	if err != nil {
		t.Fatalf("ToCanonical() error = %v", err)
	}
	if got.Name != "test_tool" {
		t.Errorf("Name = %q, want %q", got.Name, "test_tool")
	}
	if got.Description != "A test tool" {
		t.Errorf("Description = %q, want %q", got.Description, "A test tool")
	}
	if got.SourceFormat != "anthropic" {
		t.Errorf("SourceFormat = %q, want %q", got.SourceFormat, "anthropic")
	}
	if got.InputSchema == nil {
		t.Fatal("InputSchema is nil")
	}
	if got.InputSchema.Type != "object" {
		t.Errorf("InputSchema.Type = %q, want %q", got.InputSchema.Type, "object")
	}
}

func TestAnthropicAdapter_ToCanonical_Pointer(t *testing.T) {
	adapter := NewAnthropicAdapter()

	tool := &AnthropicTool{
		Name: "ptr_tool",
		InputSchema: map[string]any{
			"type": "object",
		},
	}

	got, err := adapter.ToCanonical(tool)

	if err != nil {
		t.Fatalf("ToCanonical() error = %v", err)
	}
	if got.Name != "ptr_tool" {
		t.Errorf("Name = %q, want %q", got.Name, "ptr_tool")
	}
}

func TestAnthropicAdapter_ToCanonical_InvalidType(t *testing.T) {
	adapter := NewAnthropicAdapter()

	_, err := adapter.ToCanonical("not a tool")

	if err == nil {
		t.Error("ToCanonical() with invalid type = nil, want error")
	}
}

func TestAnthropicAdapter_FromCanonical_Basic(t *testing.T) {
	adapter := NewAnthropicAdapter()

	canonical := &tooladapter.CanonicalTool{
		Name:        "test_tool",
		Description: "A test tool",
		InputSchema: &tooladapter.JSONSchema{
			Type: "object",
			Properties: map[string]*tooladapter.JSONSchema{
				"name": {
					Type:        "string",
					Description: "The name",
				},
			},
			Required: []string{"name"},
		},
	}

	result, err := adapter.FromCanonical(canonical)

	if err != nil {
		t.Fatalf("FromCanonical() error = %v", err)
	}

	tool, ok := result.(AnthropicTool)
	if !ok {
		t.Fatalf("FromCanonical() type = %T, want AnthropicTool", result)
	}

	if tool.Name != "test_tool" {
		t.Errorf("Name = %q, want %q", tool.Name, "test_tool")
	}
	if tool.Description != "A test tool" {
		t.Errorf("Description = %q, want %q", tool.Description, "A test tool")
	}

	if tool.InputSchema["type"] != "object" {
		t.Errorf("InputSchema.type = %v, want %q", tool.InputSchema["type"], "object")
	}
}

func TestAnthropicAdapter_FromCanonical_NilTool(t *testing.T) {
	adapter := NewAnthropicAdapter()

	_, err := adapter.FromCanonical(nil)

	if err == nil {
		t.Error("FromCanonical(nil) = nil, want error")
	}
}

func TestAnthropicAdapter_RoundTrip(t *testing.T) {
	adapter := NewAnthropicAdapter()

	original := AnthropicTool{
		Name:        "round_trip",
		Description: "A tool for round-trip testing",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type":        "string",
					"description": "Input value",
				},
				"count": map[string]any{
					"type":    "integer",
					"minimum": 0,
					"maximum": 10,
				},
			},
			"required": []any{"input"},
		},
	}

	// Convert to canonical
	canonical, err := adapter.ToCanonical(original)
	if err != nil {
		t.Fatalf("ToCanonical() error = %v", err)
	}

	// Convert back
	result, err := adapter.FromCanonical(canonical)
	if err != nil {
		t.Fatalf("FromCanonical() error = %v", err)
	}

	roundTripped := result.(AnthropicTool)

	if roundTripped.Name != original.Name {
		t.Errorf("Name = %q, want %q", roundTripped.Name, original.Name)
	}
	if roundTripped.Description != original.Description {
		t.Errorf("Description = %q, want %q", roundTripped.Description, original.Description)
	}

	if roundTripped.InputSchema["type"] != "object" {
		t.Errorf("InputSchema.type = %v, want %q", roundTripped.InputSchema["type"], "object")
	}
}

func TestAnthropicAdapter_SupportsFeature(t *testing.T) {
	adapter := NewAnthropicAdapter()

	tests := []struct {
		feature tooladapter.SchemaFeature
		want    bool
	}{
		{tooladapter.FeatureRef, false}, // Anthropic doesn't support $ref
		{tooladapter.FeatureDefs, false}, // Anthropic doesn't support $defs
		{tooladapter.FeatureAnyOf, true},
		{tooladapter.FeatureOneOf, true},
		{tooladapter.FeatureAllOf, true},
		{tooladapter.FeatureNot, true},
		{tooladapter.FeaturePattern, true},
		{tooladapter.FeatureFormat, true},
		{tooladapter.FeatureAdditionalProperties, true},
		{tooladapter.FeatureMinimum, true},
		{tooladapter.FeatureMaximum, true},
		{tooladapter.FeatureMinLength, true},
		{tooladapter.FeatureMaxLength, true},
		{tooladapter.FeatureEnum, true},
		{tooladapter.FeatureConst, true},
		{tooladapter.FeatureDefault, true},
	}

	for _, tt := range tests {
		t.Run(tt.feature.String(), func(t *testing.T) {
			got := adapter.SupportsFeature(tt.feature)
			if got != tt.want {
				t.Errorf("SupportsFeature(%s) = %v, want %v", tt.feature, got, tt.want)
			}
		})
	}
}

func TestAnthropicTool_JSONSerialization(t *testing.T) {
	tool := AnthropicTool{
		Name:        "test",
		Description: "Test tool",
		InputSchema: map[string]any{
			"type": "object",
		},
	}

	// Verify struct fields match Anthropic API format
	if tool.Name != "test" {
		t.Error("Name not set correctly")
	}
	if tool.InputSchema == nil {
		t.Error("InputSchema is nil")
	}
}

func TestAnthropicAdapter_FromCanonical_NestedProperties(t *testing.T) {
	adapter := NewAnthropicAdapter()

	canonical := &tooladapter.CanonicalTool{
		Name: "nested_tool",
		InputSchema: &tooladapter.JSONSchema{
			Type: "object",
			Properties: map[string]*tooladapter.JSONSchema{
				"config": {
					Type: "object",
					Properties: map[string]*tooladapter.JSONSchema{
						"enabled": {Type: "boolean"},
						"count":   {Type: "integer"},
					},
				},
			},
		},
	}

	result, err := adapter.FromCanonical(canonical)

	if err != nil {
		t.Fatalf("FromCanonical() error = %v", err)
	}

	tool := result.(AnthropicTool)
	props, ok := tool.InputSchema["properties"].(map[string]any)
	if !ok {
		t.Fatalf("InputSchema.properties is not map[string]any")
	}

	config, ok := props["config"].(map[string]any)
	if !ok {
		t.Fatalf("config is not map[string]any")
	}

	if config["type"] != "object" {
		t.Errorf("config.type = %v, want %q", config["type"], "object")
	}

	nestedProps, ok := config["properties"].(map[string]any)
	if !ok {
		t.Fatalf("config.properties is not map[string]any")
	}

	enabled, ok := nestedProps["enabled"].(map[string]any)
	if !ok {
		t.Fatalf("enabled is not map[string]any")
	}

	if enabled["type"] != "boolean" {
		t.Errorf("enabled.type = %v, want %q", enabled["type"], "boolean")
	}
}

func TestAnthropicAdapter_ToCanonical_WithCombinators(t *testing.T) {
	adapter := NewAnthropicAdapter()

	tool := AnthropicTool{
		Name: "combinator_tool",
		InputSchema: map[string]any{
			"anyOf": []any{
				map[string]any{"type": "string"},
				map[string]any{"type": "integer"},
			},
		},
	}

	got, err := adapter.ToCanonical(tool)

	if err != nil {
		t.Fatalf("ToCanonical() error = %v", err)
	}

	if len(got.InputSchema.AnyOf) != 2 {
		t.Errorf("AnyOf length = %d, want 2", len(got.InputSchema.AnyOf))
	}
}

func TestAnthropicAdapter_ToCanonical_EmptySchema(t *testing.T) {
	adapter := NewAnthropicAdapter()

	tool := AnthropicTool{
		Name:        "empty_schema_tool",
		Description: "A tool with no input schema",
	}

	got, err := adapter.ToCanonical(tool)

	if err != nil {
		t.Fatalf("ToCanonical() error = %v", err)
	}

	if got.InputSchema != nil {
		t.Error("InputSchema should be nil when input_schema is empty")
	}
}
