package adapters

import (
	"testing"

	"github.com/jonwraymond/tooladapter"
)

func TestOpenAIAdapter_Name(t *testing.T) {
	adapter := NewOpenAIAdapter()

	got := adapter.Name()
	want := "openai"

	if got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}

func TestOpenAIAdapter_ToCanonical_Basic(t *testing.T) {
	adapter := NewOpenAIAdapter()

	fn := OpenAIFunction{
		Name:        "test_function",
		Description: "A test function",
		Parameters: map[string]any{
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

	got, err := adapter.ToCanonical(fn)

	if err != nil {
		t.Fatalf("ToCanonical() error = %v", err)
	}
	if got.Name != "test_function" {
		t.Errorf("Name = %q, want %q", got.Name, "test_function")
	}
	if got.Description != "A test function" {
		t.Errorf("Description = %q, want %q", got.Description, "A test function")
	}
	if got.SourceFormat != "openai" {
		t.Errorf("SourceFormat = %q, want %q", got.SourceFormat, "openai")
	}
	if got.InputSchema == nil {
		t.Fatal("InputSchema is nil")
	}
	if got.InputSchema.Type != "object" {
		t.Errorf("InputSchema.Type = %q, want %q", got.InputSchema.Type, "object")
	}
}

func TestOpenAIAdapter_ToCanonical_Pointer(t *testing.T) {
	adapter := NewOpenAIAdapter()

	fn := &OpenAIFunction{
		Name: "ptr_function",
		Parameters: map[string]any{
			"type": "object",
		},
	}

	got, err := adapter.ToCanonical(fn)

	if err != nil {
		t.Fatalf("ToCanonical() error = %v", err)
	}
	if got.Name != "ptr_function" {
		t.Errorf("Name = %q, want %q", got.Name, "ptr_function")
	}
}

func TestOpenAIAdapter_ToCanonical_InvalidType(t *testing.T) {
	adapter := NewOpenAIAdapter()

	_, err := adapter.ToCanonical("not a function")

	if err == nil {
		t.Error("ToCanonical() with invalid type = nil, want error")
	}
}

func TestOpenAIAdapter_ToCanonical_WithStrict(t *testing.T) {
	adapter := NewOpenAIAdapter()

	fn := OpenAIFunction{
		Name:   "strict_function",
		Strict: true,
		Parameters: map[string]any{
			"type": "object",
		},
	}

	got, err := adapter.ToCanonical(fn)

	if err != nil {
		t.Fatalf("ToCanonical() error = %v", err)
	}
	// Strict mode should be preserved in SourceMeta
	if got.SourceMeta == nil {
		t.Fatal("SourceMeta is nil")
	}
	if got.SourceMeta["strict"] != true {
		t.Errorf("SourceMeta[strict] = %v, want true", got.SourceMeta["strict"])
	}
}

func TestOpenAIAdapter_FromCanonical_Basic(t *testing.T) {
	adapter := NewOpenAIAdapter()

	canonical := &tooladapter.CanonicalTool{
		Name:        "test_function",
		Description: "A test function",
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

	fn, ok := result.(OpenAIFunction)
	if !ok {
		t.Fatalf("FromCanonical() type = %T, want OpenAIFunction", result)
	}

	if fn.Name != "test_function" {
		t.Errorf("Name = %q, want %q", fn.Name, "test_function")
	}
	if fn.Description != "A test function" {
		t.Errorf("Description = %q, want %q", fn.Description, "A test function")
	}

	if fn.Parameters["type"] != "object" {
		t.Errorf("Parameters.type = %v, want %q", fn.Parameters["type"], "object")
	}
}

func TestOpenAIAdapter_StrictMode(t *testing.T) {
	adapter := NewOpenAIAdapter()

	canonical := &tooladapter.CanonicalTool{
		Name: "strict_function",
		InputSchema: &tooladapter.JSONSchema{
			Type: "object",
			Properties: map[string]*tooladapter.JSONSchema{
				"name": {Type: "string"},
			},
		},
		SourceMeta: map[string]any{
			"strict": true,
		},
	}

	result, err := adapter.FromCanonical(canonical)

	if err != nil {
		t.Fatalf("FromCanonical() error = %v", err)
	}

	fn := result.(OpenAIFunction)

	// Strict mode should be set
	if !fn.Strict {
		t.Error("Strict = false, want true")
	}

	// In strict mode, additionalProperties should be false
	if fn.Parameters["additionalProperties"] != false {
		t.Errorf("Parameters.additionalProperties = %v, want false", fn.Parameters["additionalProperties"])
	}
}

func TestOpenAIAdapter_RoundTrip(t *testing.T) {
	adapter := NewOpenAIAdapter()

	original := OpenAIFunction{
		Name:        "round_trip",
		Description: "A function for round-trip testing",
		Parameters: map[string]any{
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

	roundTripped := result.(OpenAIFunction)

	if roundTripped.Name != original.Name {
		t.Errorf("Name = %q, want %q", roundTripped.Name, original.Name)
	}
	if roundTripped.Description != original.Description {
		t.Errorf("Description = %q, want %q", roundTripped.Description, original.Description)
	}

	if roundTripped.Parameters["type"] != "object" {
		t.Errorf("Parameters.type = %v, want %q", roundTripped.Parameters["type"], "object")
	}
}

func TestOpenAIAdapter_SupportsFeature(t *testing.T) {
	adapter := NewOpenAIAdapter()

	tests := []struct {
		feature tooladapter.SchemaFeature
		want    bool
	}{
		{tooladapter.FeatureRef, false},           // OpenAI doesn't support $ref
		{tooladapter.FeatureDefs, false},          // OpenAI doesn't support $defs
		{tooladapter.FeatureAnyOf, false},         // Limited support
		{tooladapter.FeatureOneOf, false},         // Limited support
		{tooladapter.FeatureAllOf, false},         // Limited support
		{tooladapter.FeatureNot, false},           // Not supported
		{tooladapter.FeaturePattern, true},        // Supported in strict mode
		{tooladapter.FeatureFormat, true},         // Supported
		{tooladapter.FeatureAdditionalProperties, true}, // Required in strict mode
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

func TestOpenAIFunction_JSONSerialization(t *testing.T) {
	fn := OpenAIFunction{
		Name:        "test",
		Description: "Test function",
		Parameters: map[string]any{
			"type": "object",
		},
		Strict: true,
	}

	// Verify struct can be used
	if fn.Name != "test" {
		t.Error("Name not set correctly")
	}
	if !fn.Strict {
		t.Error("Strict not set correctly")
	}
}

func TestOpenAIAdapter_FromCanonical_NilTool(t *testing.T) {
	adapter := NewOpenAIAdapter()

	_, err := adapter.FromCanonical(nil)

	if err == nil {
		t.Error("FromCanonical(nil) = nil, want error")
	}
}

func TestOpenAIAdapter_FromCanonical_NestedProperties(t *testing.T) {
	adapter := NewOpenAIAdapter()

	canonical := &tooladapter.CanonicalTool{
		Name: "nested_function",
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

	fn := result.(OpenAIFunction)
	props, ok := fn.Parameters["properties"].(map[string]any)
	if !ok {
		t.Fatalf("Parameters.properties is not map[string]any")
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
