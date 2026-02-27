package entity_test

import (
	"testing"

	"github.com/jjmrocha/knowledge-mcp/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestParseContent_ValidContent(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedMetadata string
		expectedBody     string
	}{
		{
			name:             "standard frontmatter with body",
			input:            "---\nentity: concept\nschema: 1\n---\nThis is the body.\n",
			expectedMetadata: "entity: concept\nschema: 1",
			expectedBody:     "This is the body.\n",
		},
		{
			name:             "empty metadata",
			input:            "---\n\n---\nBody text.",
			expectedMetadata: "",
			expectedBody:     "Body text.",
		},
		{
			name:             "empty body",
			input:            "---\nentity: tag\n---\n",
			expectedMetadata: "entity: tag",
			expectedBody:     "",
		},
		{
			name:             "empty metadata and empty body",
			input:            "---\n\n---\n",
			expectedMetadata: "",
			expectedBody:     "",
		},
		{
			name:             "multi-line metadata and multi-line body",
			input:            "---\nentity: concept\nschema: 1\nuri: scio://contexts/ecommerce/domains/business-rules/concepts/discount\nname: Discount\n---\nFirst paragraph.\n\nSecond paragraph.\n",
			expectedMetadata: "entity: concept\nschema: 1\nuri: scio://contexts/ecommerce/domains/business-rules/concepts/discount\nname: Discount",
			expectedBody:     "First paragraph.\n\nSecond paragraph.\n",
		},
		{
			name:             "dashes inside body do not break parsing",
			input:            "---\nentity: concept\n---\nSome text.\n---\nMore text after dashes.\n",
			expectedMetadata: "entity: concept",
			expectedBody:     "Some text.\n---\nMore text after dashes.\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := entity.ParseContent(tc.input)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tc.expectedMetadata, result.Metadata)
			assert.Equal(t, tc.expectedBody, result.Body)
		})
	}
}

func TestParseContent_InvalidContent(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty string",
			input: "",
		},
		{
			name:  "no frontmatter delimiter",
			input: "just plain text",
		},
		{
			name:  "opening delimiter without newline",
			input: "---entity: tag---",
		},
		{
			name:  "only opening delimiter",
			input: "---\nentity: tag",
		},
		{
			name:  "opening delimiter present but no closing delimiter",
			input: "---\nentity: tag\n",
		},
		{
			name:  "closing delimiter without preceding newline",
			input: "---\nentity: tag---\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := entity.ParseContent(tc.input)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestEncode_ValidContent(t *testing.T) {
	tests := []struct {
		name     string
		input    entity.EntityContent
		expected string
	}{
		{
			name: "metadata and body",
			input: entity.EntityContent{
				Metadata: "entity: concept\nschema: 1",
				Body:     "Body text.\n",
			},
			expected: "---\nentity: concept\nschema: 1\n---\nBody text.\n",
		},
		{
			name: "empty metadata",
			input: entity.EntityContent{
				Metadata: "",
				Body:     "Body only.\n",
			},
			expected: "---\n\n---\nBody only.\n",
		},
		{
			name: "empty body",
			input: entity.EntityContent{
				Metadata: "entity: tag",
				Body:     "",
			},
			expected: "---\nentity: tag\n---\n",
		},
		{
			name: "empty metadata and empty body",
			input: entity.EntityContent{
				Metadata: "",
				Body:     "",
			},
			expected: "---\n\n---\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := entity.Encode(&tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestParseContent_RoundTrip(t *testing.T) {
	// given
	original := entity.EntityContent{
		Metadata: `entity: concept
		schema: 1
		uri: scio://contexts/ecommerce/domains/rules/concepts/discount
		name: Discount Calculation`,
		Body: `Business logic for calculating tiered discounts based on order total.`,
	}
	// when
	encoded := entity.Encode(&original)

	parsed, err := entity.ParseContent(encoded)
	// then
	assert.NoError(t, err)
	assert.NotNil(t, parsed)
	assert.Equal(t, original.Metadata, parsed.Metadata)
	assert.Equal(t, original.Body, parsed.Body)
}
