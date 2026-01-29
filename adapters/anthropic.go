package adapters

import (
	"errors"

	"github.com/jonwraymond/tooladapter"
)

// AnthropicTool represents an Anthropic tool definition.
// This is a self-contained type that doesn't depend on external SDK.
type AnthropicTool struct {
	// Name is the tool identifier
	Name string `json:"name"`

	// Description explains what the tool does
	Description string `json:"description,omitempty"`

	// InputSchema is the JSON Schema for tool input (Anthropic uses input_schema)
	InputSchema map[string]any `json:"input_schema"`
}

// AnthropicAdapter converts between Anthropic tool format and canonical format.
type AnthropicAdapter struct{}

// NewAnthropicAdapter creates a new Anthropic adapter.
func NewAnthropicAdapter() *AnthropicAdapter {
	return &AnthropicAdapter{}
}

// Name returns the adapter identifier.
func (a *AnthropicAdapter) Name() string {
	return "anthropic"
}

// ToCanonical converts an Anthropic tool to canonical format.
// Accepts AnthropicTool or *AnthropicTool.
func (a *AnthropicAdapter) ToCanonical(raw any) (*tooladapter.CanonicalTool, error) {
	var tool AnthropicTool

	switch v := raw.(type) {
	case AnthropicTool:
		tool = v
	case *AnthropicTool:
		if v == nil {
			return nil, errors.New("nil AnthropicTool pointer")
		}
		tool = *v
	default:
		return nil, errors.New("expected AnthropicTool or *AnthropicTool")
	}

	canonical := &tooladapter.CanonicalTool{
		Name:         tool.Name,
		Description:  tool.Description,
		SourceFormat: "anthropic",
		SourceMeta:   make(map[string]any),
	}

	// Convert input schema
	if tool.InputSchema != nil {
		schema, err := mapToJSONSchema(tool.InputSchema)
		if err != nil {
			return nil, err
		}
		canonical.InputSchema = schema
	}

	return canonical, nil
}

// FromCanonical converts a canonical tool to Anthropic format.
func (a *AnthropicAdapter) FromCanonical(tool *tooladapter.CanonicalTool) (any, error) {
	if tool == nil {
		return nil, errors.New("nil CanonicalTool")
	}

	anthropicTool := AnthropicTool{
		Name:        tool.Name,
		Description: tool.Description,
	}

	// Convert input schema to input_schema map
	if tool.InputSchema != nil {
		anthropicTool.InputSchema = tool.InputSchema.ToMap()
	}

	return anthropicTool, nil
}

// SupportsFeature returns whether this adapter supports a schema feature.
// Anthropic supports most JSON Schema features except $ref.
func (a *AnthropicAdapter) SupportsFeature(feature tooladapter.SchemaFeature) bool {
	switch feature {
	case tooladapter.FeatureRef:
		return false // Anthropic doesn't support $ref
	case tooladapter.FeatureDefs:
		return false // Anthropic doesn't support $defs
	default:
		return true // Other features are generally supported
	}
}
