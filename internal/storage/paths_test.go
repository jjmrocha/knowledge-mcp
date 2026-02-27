package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jjmrocha/knowledge-mcp/internal/storage"
	"github.com/jjmrocha/knowledge-mcp/internal/uri"
)

const rootDir = "/data/knowledge"

func TestFileName(t *testing.T) {
	tests := []struct {
		name     string
		rawURI   string
		expected string
	}{
		{
			name:     "global tag",
			rawURI:   "scio://tags/business-rule",
			expected: "/data/knowledge/tags/business-rule.md",
		},
		{
			name:     "context-scoped tag",
			rawURI:   "scio://contexts/ecommerce/tags/local-tag",
			expected: "/data/knowledge/contexts/ecommerce/tags/local-tag.md",
		},
		{
			name:     "global relation",
			rawURI:   "scio://relations/implements",
			expected: "/data/knowledge/relations/implements.md",
		},
		{
			name:     "context-scoped relation",
			rawURI:   "scio://contexts/ecommerce/relations/owns",
			expected: "/data/knowledge/contexts/ecommerce/relations/owns.md",
		},
		{
			name:     "context",
			rawURI:   "scio://contexts/ecommerce",
			expected: "/data/knowledge/contexts/ecommerce/context.md",
		},
		{
			name:     "domain",
			rawURI:   "scio://contexts/ecommerce/domains/business-rules",
			expected: "/data/knowledge/contexts/ecommerce/domains/business-rules/domain.md",
		},
		{
			name:     "concept",
			rawURI:   "scio://contexts/ecommerce/domains/business-rules/concepts/discount-calculation",
			expected: "/data/knowledge/contexts/ecommerce/domains/business-rules/discount-calculation.md",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, err := uri.Parse(tc.rawURI)
			assert.NoError(t, err)

			result := storage.FileName(rootDir, u)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFileDir(t *testing.T) {
	tests := []struct {
		name     string
		rawURI   string
		expected string
	}{
		{
			name:     "global tag",
			rawURI:   "scio://tags/business-rule",
			expected: "/data/knowledge/tags",
		},
		{
			name:     "context-scoped tag",
			rawURI:   "scio://contexts/ecommerce/tags/local-tag",
			expected: "/data/knowledge/contexts/ecommerce/tags",
		},
		{
			name:     "global relation",
			rawURI:   "scio://relations/implements",
			expected: "/data/knowledge/relations",
		},
		{
			name:     "context-scoped relation",
			rawURI:   "scio://contexts/ecommerce/relations/owns",
			expected: "/data/knowledge/contexts/ecommerce/relations",
		},
		{
			name:     "context",
			rawURI:   "scio://contexts/ecommerce",
			expected: "/data/knowledge/contexts/ecommerce",
		},
		{
			name:     "domain",
			rawURI:   "scio://contexts/ecommerce/domains/business-rules",
			expected: "/data/knowledge/contexts/ecommerce/domains/business-rules",
		},
		{
			name:     "concept",
			rawURI:   "scio://contexts/ecommerce/domains/business-rules/concepts/discount-calculation",
			expected: "/data/knowledge/contexts/ecommerce/domains/business-rules",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, err := uri.Parse(tc.rawURI)
			assert.NoError(t, err)

			result := storage.FileDir(rootDir, u)
			assert.Equal(t, tc.expected, result)
		})
	}
}
