package adapters

import (
	"errors"

	"github.com/jonwraymond/tooladapter"
)

// OpenAIFunction represents an OpenAI function/tool definition.
// This is a self-contained type that doesn't depend on external SDK.
type OpenAIFunction struct {
	// Name is the function identifier
	Name string `json:"name"`

	// Description explains what the function does
	Description string `json:"description,omitempty"`

	// Parameters is the JSON Schema for function arguments
	Parameters map[string]any `json:"parameters"`

	// Strict enables strict mode (additionalProperties=false enforced)
	Strict bool `json:"strict,omitempty"`
}

// OpenAIAdapter converts between OpenAI function format and canonical format.
type OpenAIAdapter struct{}

// NewOpenAIAdapter creates a new OpenAI adapter.
func NewOpenAIAdapter() *OpenAIAdapter {
	return &OpenAIAdapter{}
}

// Name returns the adapter identifier.
func (a *OpenAIAdapter) Name() string {
	return "openai"
}

// ToCanonical converts an OpenAI function to canonical format.
// Accepts OpenAIFunction or *OpenAIFunction.
func (a *OpenAIAdapter) ToCanonical(raw any) (*tooladapter.CanonicalTool, error) {
	var fn OpenAIFunction

	switch v := raw.(type) {
	case OpenAIFunction:
		fn = v
	case *OpenAIFunction:
		if v == nil {
			return nil, errors.New("nil OpenAIFunction pointer")
		}
		fn = *v
	default:
		return nil, errors.New("expected OpenAIFunction or *OpenAIFunction")
	}

	canonical := &tooladapter.CanonicalTool{
		Name:         fn.Name,
		Description:  fn.Description,
		SourceFormat: "openai",
		SourceMeta:   make(map[string]any),
	}

	// Store OpenAI-specific fields in SourceMeta for round-trip
	if fn.Strict {
		canonical.SourceMeta["strict"] = true
	}

	// Convert parameters schema
	if fn.Parameters != nil {
		schema, err := mapToJSONSchema(fn.Parameters)
		if err != nil {
			return nil, err
		}
		canonical.InputSchema = schema
	}

	return canonical, nil
}

// FromCanonical converts a canonical tool to OpenAI format.
func (a *OpenAIAdapter) FromCanonical(tool *tooladapter.CanonicalTool) (any, error) {
	if tool == nil {
		return nil, errors.New("nil CanonicalTool")
	}

	fn := OpenAIFunction{
		Name:        tool.Name,
		Description: tool.Description,
	}

	// Restore OpenAI-specific fields from SourceMeta
	if tool.SourceMeta != nil {
		if strict, ok := tool.SourceMeta["strict"].(bool); ok && strict {
			fn.Strict = true
		}
	}

	// Convert input schema to parameters map
	if tool.InputSchema != nil {
		params := tool.InputSchema.ToMap()

		// In strict mode, enforce additionalProperties=false at root
		if fn.Strict {
			params["additionalProperties"] = false
		}

		fn.Parameters = params
	}

	return fn, nil
}

// SupportsFeature returns whether this adapter supports a schema feature.
// OpenAI has limited schema support compared to full JSON Schema.
func (a *OpenAIAdapter) SupportsFeature(feature tooladapter.SchemaFeature) bool {
	switch feature {
	case tooladapter.FeatureRef:
		return false // OpenAI doesn't support $ref
	case tooladapter.FeatureDefs:
		return false // OpenAI doesn't support $defs
	case tooladapter.FeatureAnyOf:
		return false // Limited/no support
	case tooladapter.FeatureOneOf:
		return false // Limited/no support
	case tooladapter.FeatureAllOf:
		return false // Limited/no support
	case tooladapter.FeatureNot:
		return false // Not supported
	default:
		return true // Other features are generally supported
	}
}
