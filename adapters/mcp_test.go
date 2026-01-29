package adapters

import (
	"testing"

	"github.com/jonwraymond/tooladapter"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPAdapter_Name(t *testing.T) {
	adapter := NewMCPAdapter()

	got := adapter.Name()
	want := "mcp"

	if got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}

func TestMCPAdapter_ToCanonical_Basic(t *testing.T) {
	adapter := NewMCPAdapter()

	mcpTool := mcp.Tool{
		Name:        "test-tool",
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

	got, err := adapter.ToCanonical(mcpTool)

	if err != nil {
		t.Fatalf("ToCanonical() error = %v", err)
	}
	if got.Name != "test-tool" {
		t.Errorf("Name = %q, want %q", got.Name, "test-tool")
	}
	if got.Description != "A test tool" {
		t.Errorf("Description = %q, want %q", got.Description, "A test tool")
	}
	if got.SourceFormat != "mcp" {
		t.Errorf("SourceFormat = %q, want %q", got.SourceFormat, "mcp")
	}
	if got.InputSchema == nil {
		t.Fatal("InputSchema is nil")
	}
	if got.InputSchema.Type != "object" {
		t.Errorf("InputSchema.Type = %q, want %q", got.InputSchema.Type, "object")
	}
	if len(got.InputSchema.Properties) != 1 {
		t.Errorf("InputSchema.Properties length = %d, want 1", len(got.InputSchema.Properties))
	}
}

func TestMCPAdapter_ToCanonical_Pointer(t *testing.T) {
	adapter := NewMCPAdapter()

	mcpTool := &mcp.Tool{
		Name:        "ptr-tool",
		Description: "A pointer tool",
		InputSchema: map[string]any{
			"type": "object",
		},
	}

	got, err := adapter.ToCanonical(mcpTool)

	if err != nil {
		t.Fatalf("ToCanonical() error = %v", err)
	}
	if got.Name != "ptr-tool" {
		t.Errorf("Name = %q, want %q", got.Name, "ptr-tool")
	}
}

func TestMCPAdapter_ToCanonical_InvalidType(t *testing.T) {
	adapter := NewMCPAdapter()

	_, err := adapter.ToCanonical("not a tool")

	if err == nil {
		t.Error("ToCanonical() with invalid type = nil, want error")
	}
}

func TestMCPAdapter_ToCanonical_WithTitle(t *testing.T) {
	adapter := NewMCPAdapter()

	mcpTool := mcp.Tool{
		Name:        "tool-name",
		Title:       "Tool Title",
		Description: "Description",
		InputSchema: map[string]any{"type": "object"},
	}

	got, err := adapter.ToCanonical(mcpTool)

	if err != nil {
		t.Fatalf("ToCanonical() error = %v", err)
	}
	// Title should be preserved in SourceMeta
	if got.SourceMeta == nil {
		t.Fatal("SourceMeta is nil")
	}
	if got.SourceMeta["title"] != "Tool Title" {
		t.Errorf("SourceMeta[title] = %v, want %q", got.SourceMeta["title"], "Tool Title")
	}
}

func TestMCPAdapter_ToCanonical_WithOutputSchema(t *testing.T) {
	adapter := NewMCPAdapter()

	mcpTool := mcp.Tool{
		Name: "output-tool",
		InputSchema: map[string]any{
			"type": "object",
		},
		OutputSchema: map[string]any{
			"type": "string",
		},
	}

	got, err := adapter.ToCanonical(mcpTool)

	if err != nil {
		t.Fatalf("ToCanonical() error = %v", err)
	}
	if got.OutputSchema == nil {
		t.Fatal("OutputSchema is nil")
	}
	if got.OutputSchema.Type != "string" {
		t.Errorf("OutputSchema.Type = %q, want %q", got.OutputSchema.Type, "string")
	}
}

func TestMCPAdapter_FromCanonical_Basic(t *testing.T) {
	adapter := NewMCPAdapter()

	canonical := &tooladapter.CanonicalTool{
		Name:        "test-tool",
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

	mcpTool, ok := result.(mcp.Tool)
	if !ok {
		t.Fatalf("FromCanonical() type = %T, want mcp.Tool", result)
	}

	if mcpTool.Name != "test-tool" {
		t.Errorf("Name = %q, want %q", mcpTool.Name, "test-tool")
	}
	if mcpTool.Description != "A test tool" {
		t.Errorf("Description = %q, want %q", mcpTool.Description, "A test tool")
	}

	schema, ok := mcpTool.InputSchema.(map[string]any)
	if !ok {
		t.Fatalf("InputSchema type = %T, want map[string]any", mcpTool.InputSchema)
	}
	if schema["type"] != "object" {
		t.Errorf("InputSchema.type = %v, want %q", schema["type"], "object")
	}
}

func TestMCPAdapter_FromCanonical_WithOutputSchema(t *testing.T) {
	adapter := NewMCPAdapter()

	canonical := &tooladapter.CanonicalTool{
		Name: "output-tool",
		InputSchema: &tooladapter.JSONSchema{
			Type: "object",
		},
		OutputSchema: &tooladapter.JSONSchema{
			Type: "string",
		},
	}

	result, err := adapter.FromCanonical(canonical)

	if err != nil {
		t.Fatalf("FromCanonical() error = %v", err)
	}

	mcpTool := result.(mcp.Tool)
	if mcpTool.OutputSchema == nil {
		t.Fatal("OutputSchema is nil")
	}

	schema, ok := mcpTool.OutputSchema.(map[string]any)
	if !ok {
		t.Fatalf("OutputSchema type = %T, want map[string]any", mcpTool.OutputSchema)
	}
	if schema["type"] != "string" {
		t.Errorf("OutputSchema.type = %v, want %q", schema["type"], "string")
	}
}

func TestMCPAdapter_FromCanonical_WithTitle(t *testing.T) {
	adapter := NewMCPAdapter()

	canonical := &tooladapter.CanonicalTool{
		Name: "titled-tool",
		InputSchema: &tooladapter.JSONSchema{
			Type: "object",
		},
		SourceMeta: map[string]any{
			"title": "Tool Title",
		},
	}

	result, err := adapter.FromCanonical(canonical)

	if err != nil {
		t.Fatalf("FromCanonical() error = %v", err)
	}

	mcpTool := result.(mcp.Tool)
	if mcpTool.Title != "Tool Title" {
		t.Errorf("Title = %q, want %q", mcpTool.Title, "Tool Title")
	}
}

func TestMCPAdapter_RoundTrip(t *testing.T) {
	adapter := NewMCPAdapter()

	original := mcp.Tool{
		Name:        "round-trip",
		Title:       "Round Trip Tool",
		Description: "A tool for round-trip testing",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type":        "string",
					"description": "Input value",
					"minLength":   1,
					"maxLength":   100,
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

	roundTripped := result.(mcp.Tool)

	// Verify key fields
	if roundTripped.Name != original.Name {
		t.Errorf("Name = %q, want %q", roundTripped.Name, original.Name)
	}
	if roundTripped.Title != original.Title {
		t.Errorf("Title = %q, want %q", roundTripped.Title, original.Title)
	}
	if roundTripped.Description != original.Description {
		t.Errorf("Description = %q, want %q", roundTripped.Description, original.Description)
	}

	// Verify schema structure
	schema, ok := roundTripped.InputSchema.(map[string]any)
	if !ok {
		t.Fatalf("InputSchema type = %T, want map[string]any", roundTripped.InputSchema)
	}
	if schema["type"] != "object" {
		t.Errorf("InputSchema.type = %v, want %q", schema["type"], "object")
	}
}

func TestMCPAdapter_SupportsFeature(t *testing.T) {
	adapter := NewMCPAdapter()

	// MCP should support all features
	for _, feature := range tooladapter.AllFeatures() {
		if !adapter.SupportsFeature(feature) {
			t.Errorf("SupportsFeature(%s) = false, want true", feature)
		}
	}
}

func TestMCPAdapter_ToCanonical_ComplexSchema(t *testing.T) {
	adapter := NewMCPAdapter()

	mcpTool := mcp.Tool{
		Name: "complex-tool",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"items": map[string]any{
					"type": "array",
					"items": map[string]any{
						"type": "string",
					},
				},
				"options": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"enabled": map[string]any{
							"type":    "boolean",
							"default": true,
						},
					},
				},
			},
			"$defs": map[string]any{
				"Status": map[string]any{
					"type": "string",
					"enum": []any{"active", "inactive"},
				},
			},
		},
	}

	got, err := adapter.ToCanonical(mcpTool)

	if err != nil {
		t.Fatalf("ToCanonical() error = %v", err)
	}

	// Verify nested array
	if got.InputSchema.Properties["items"] == nil {
		t.Fatal("items property is nil")
	}
	if got.InputSchema.Properties["items"].Type != "array" {
		t.Error("items.type != array")
	}
	if got.InputSchema.Properties["items"].Items == nil {
		t.Error("items.items is nil")
	}

	// Verify $defs
	if len(got.InputSchema.Defs) != 1 {
		t.Errorf("Defs length = %d, want 1", len(got.InputSchema.Defs))
	}
}

func TestMCPAdapter_ToCanonical_WithRef(t *testing.T) {
	adapter := NewMCPAdapter()

	mcpTool := mcp.Tool{
		Name: "ref-tool",
		InputSchema: map[string]any{
			"$ref": "#/$defs/Person",
			"$defs": map[string]any{
				"Person": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"name": map[string]any{"type": "string"},
					},
				},
			},
		},
	}

	got, err := adapter.ToCanonical(mcpTool)

	if err != nil {
		t.Fatalf("ToCanonical() error = %v", err)
	}

	if got.InputSchema.Ref != "#/$defs/Person" {
		t.Errorf("Ref = %q, want %q", got.InputSchema.Ref, "#/$defs/Person")
	}
}

func TestMCPAdapter_ToCanonical_WithCombinators(t *testing.T) {
	adapter := NewMCPAdapter()

	mcpTool := mcp.Tool{
		Name: "combinator-tool",
		InputSchema: map[string]any{
			"anyOf": []any{
				map[string]any{"type": "string"},
				map[string]any{"type": "integer"},
			},
		},
	}

	got, err := adapter.ToCanonical(mcpTool)

	if err != nil {
		t.Fatalf("ToCanonical() error = %v", err)
	}

	if len(got.InputSchema.AnyOf) != 2 {
		t.Errorf("AnyOf length = %d, want 2", len(got.InputSchema.AnyOf))
	}
}
