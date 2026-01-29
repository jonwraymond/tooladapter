package adapters

import (
	"errors"

	"github.com/jonwraymond/tooladapter"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPAdapter converts between MCP tool format and canonical format.
type MCPAdapter struct{}

// NewMCPAdapter creates a new MCP adapter.
func NewMCPAdapter() *MCPAdapter {
	return &MCPAdapter{}
}

// Name returns the adapter identifier.
func (a *MCPAdapter) Name() string {
	return "mcp"
}

// ToCanonical converts an MCP tool to canonical format.
// Accepts mcp.Tool or *mcp.Tool.
func (a *MCPAdapter) ToCanonical(raw any) (*tooladapter.CanonicalTool, error) {
	var tool mcp.Tool

	switch v := raw.(type) {
	case mcp.Tool:
		tool = v
	case *mcp.Tool:
		if v == nil {
			return nil, errors.New("nil mcp.Tool pointer")
		}
		tool = *v
	default:
		return nil, errors.New("expected mcp.Tool or *mcp.Tool")
	}

	canonical := &tooladapter.CanonicalTool{
		Name:         tool.Name,
		Description:  tool.Description,
		SourceFormat: "mcp",
		SourceMeta:   make(map[string]any),
	}

	// Store MCP-specific fields in SourceMeta for round-trip
	if tool.Title != "" {
		canonical.SourceMeta["title"] = tool.Title
	}

	// Convert input schema
	if tool.InputSchema != nil {
		schema, err := mapToJSONSchema(tool.InputSchema)
		if err != nil {
			return nil, err
		}
		canonical.InputSchema = schema
	}

	// Convert output schema
	if tool.OutputSchema != nil {
		schema, err := mapToJSONSchema(tool.OutputSchema)
		if err != nil {
			return nil, err
		}
		canonical.OutputSchema = schema
	}

	return canonical, nil
}

// FromCanonical converts a canonical tool to MCP format.
func (a *MCPAdapter) FromCanonical(tool *tooladapter.CanonicalTool) (any, error) {
	if tool == nil {
		return nil, errors.New("nil CanonicalTool")
	}

	mcpTool := mcp.Tool{
		Name:        tool.Name,
		Description: tool.Description,
	}

	// Restore MCP-specific fields from SourceMeta
	if tool.SourceMeta != nil {
		if title, ok := tool.SourceMeta["title"].(string); ok {
			mcpTool.Title = title
		}
	}

	// Convert input schema to map
	if tool.InputSchema != nil {
		mcpTool.InputSchema = tool.InputSchema.ToMap()
	}

	// Convert output schema to map
	if tool.OutputSchema != nil {
		mcpTool.OutputSchema = tool.OutputSchema.ToMap()
	}

	return mcpTool, nil
}

// SupportsFeature returns true for all features since MCP supports the full
// JSON Schema specification.
func (a *MCPAdapter) SupportsFeature(feature tooladapter.SchemaFeature) bool {
	return true
}

// mapToJSONSchema converts a map[string]any schema to JSONSchema.
func mapToJSONSchema(raw any) (*tooladapter.JSONSchema, error) {
	m, ok := raw.(map[string]any)
	if !ok {
		return nil, errors.New("schema is not a map[string]any")
	}

	schema := &tooladapter.JSONSchema{}

	// Type
	if v, ok := m["type"].(string); ok {
		schema.Type = v
	}

	// Description
	if v, ok := m["description"].(string); ok {
		schema.Description = v
	}

	// Pattern
	if v, ok := m["pattern"].(string); ok {
		schema.Pattern = v
	}

	// Format
	if v, ok := m["format"].(string); ok {
		schema.Format = v
	}

	// $ref
	if v, ok := m["$ref"].(string); ok {
		schema.Ref = v
	}

	// Minimum
	if v, ok := m["minimum"].(float64); ok {
		schema.Minimum = &v
	} else if v, ok := m["minimum"].(int); ok {
		f := float64(v)
		schema.Minimum = &f
	}

	// Maximum
	if v, ok := m["maximum"].(float64); ok {
		schema.Maximum = &v
	} else if v, ok := m["maximum"].(int); ok {
		f := float64(v)
		schema.Maximum = &f
	}

	// MinLength
	if v, ok := m["minLength"].(float64); ok {
		i := int(v)
		schema.MinLength = &i
	} else if v, ok := m["minLength"].(int); ok {
		schema.MinLength = &v
	}

	// MaxLength
	if v, ok := m["maxLength"].(float64); ok {
		i := int(v)
		schema.MaxLength = &i
	} else if v, ok := m["maxLength"].(int); ok {
		schema.MaxLength = &v
	}

	// Const
	if v, ok := m["const"]; ok {
		schema.Const = v
	}

	// Default
	if v, ok := m["default"]; ok {
		schema.Default = v
	}

	// AdditionalProperties
	if v, ok := m["additionalProperties"].(bool); ok {
		schema.AdditionalProperties = &v
	}

	// Enum
	if v, ok := m["enum"].([]any); ok {
		schema.Enum = v
	}

	// Required
	if v, ok := m["required"].([]any); ok {
		schema.Required = make([]string, 0, len(v))
		for _, r := range v {
			if s, ok := r.(string); ok {
				schema.Required = append(schema.Required, s)
			}
		}
	}

	// Properties
	if v, ok := m["properties"].(map[string]any); ok {
		schema.Properties = make(map[string]*tooladapter.JSONSchema, len(v))
		for name, prop := range v {
			propSchema, err := mapToJSONSchema(prop)
			if err != nil {
				return nil, err
			}
			schema.Properties[name] = propSchema
		}
	}

	// Items
	if v, ok := m["items"]; ok {
		itemSchema, err := mapToJSONSchema(v)
		if err != nil {
			return nil, err
		}
		schema.Items = itemSchema
	}

	// $defs
	if v, ok := m["$defs"].(map[string]any); ok {
		schema.Defs = make(map[string]*tooladapter.JSONSchema, len(v))
		for name, def := range v {
			defSchema, err := mapToJSONSchema(def)
			if err != nil {
				return nil, err
			}
			schema.Defs[name] = defSchema
		}
	}

	// anyOf
	if v, ok := m["anyOf"].([]any); ok {
		schema.AnyOf = make([]*tooladapter.JSONSchema, 0, len(v))
		for _, item := range v {
			itemSchema, err := mapToJSONSchema(item)
			if err != nil {
				return nil, err
			}
			schema.AnyOf = append(schema.AnyOf, itemSchema)
		}
	}

	// oneOf
	if v, ok := m["oneOf"].([]any); ok {
		schema.OneOf = make([]*tooladapter.JSONSchema, 0, len(v))
		for _, item := range v {
			itemSchema, err := mapToJSONSchema(item)
			if err != nil {
				return nil, err
			}
			schema.OneOf = append(schema.OneOf, itemSchema)
		}
	}

	// allOf
	if v, ok := m["allOf"].([]any); ok {
		schema.AllOf = make([]*tooladapter.JSONSchema, 0, len(v))
		for _, item := range v {
			itemSchema, err := mapToJSONSchema(item)
			if err != nil {
				return nil, err
			}
			schema.AllOf = append(schema.AllOf, itemSchema)
		}
	}

	// not
	if v, ok := m["not"]; ok {
		notSchema, err := mapToJSONSchema(v)
		if err != nil {
			return nil, err
		}
		schema.Not = notSchema
	}

	return schema, nil
}
