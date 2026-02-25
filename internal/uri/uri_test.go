package uri_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jjmrocha/knowledge-mcp/internal/uri"
)

func TestParse_ValidPatterns(t *testing.T) {
	tests := []struct {
		name           string
		uri            string
		expectedEntity uri.EntityType
		expectedCtx    *string
		expectedDomain *string
		expectedSlug   string
	}{
		{
			name:           "global tag",
			uri:            "scio://tags/business-rule",
			expectedEntity: uri.EntityTypeTag,
			expectedSlug:   "business-rule",
		},
		{
			name:           "context-scoped tag",
			uri:            "scio://contexts/ecommerce/tags/local-tag",
			expectedEntity: uri.EntityTypeTag,
			expectedCtx:    new("ecommerce"),
			expectedSlug:   "local-tag",
		},
		{
			name:           "global relation",
			uri:            "scio://relations/implements",
			expectedEntity: uri.EntityTypeRelation,
			expectedSlug:   "implements",
		},
		{
			name:           "context-scoped relation",
			uri:            "scio://contexts/ecommerce/relations/owns",
			expectedEntity: uri.EntityTypeRelation,
			expectedCtx:    new("ecommerce"),
			expectedSlug:   "owns",
		},
		{
			name:           "context",
			uri:            "scio://contexts/ecommerce",
			expectedEntity: uri.EntityTypeContext,
			expectedSlug:   "ecommerce",
		},
		{
			name:           "domain",
			uri:            "scio://contexts/ecommerce/domains/business-rules",
			expectedEntity: uri.EntityTypeDomain,
			expectedCtx:    new("ecommerce"),
			expectedSlug:   "business-rules",
		},
		{
			name:           "concept",
			uri:            "scio://contexts/ecommerce/domains/business-rules/concepts/discount-calculation",
			expectedEntity: uri.EntityTypeConcept,
			expectedCtx:    new("ecommerce"),
			expectedDomain: new("business-rules"),
			expectedSlug:   "discount-calculation",
		},
		{
			name:           "single-char slug",
			uri:            "scio://tags/a",
			expectedEntity: uri.EntityTypeTag,
			expectedSlug:   "a",
		},
		{
			name:           "slug with digits",
			uri:            "scio://tags/rule-v2",
			expectedEntity: uri.EntityTypeTag,
			expectedSlug:   "rule-v2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := uri.Parse(tc.uri)
			assert.NoError(t, err)
			assert.NotNil(t, result)

			assert.Equal(t, tc.uri, result.Raw)
			assert.Equal(t, tc.expectedEntity, result.Entity)
			assert.Equal(t, tc.expectedSlug, result.Slug)

			if tc.expectedCtx == nil {
				assert.Nil(t, result.Context)
			} else {
				assert.NotNil(t, result.Context)
				assert.Equal(t, *tc.expectedCtx, *result.Context)
			}

			if tc.expectedDomain == nil {
				assert.Nil(t, result.Domain)
			} else {
				assert.NotNil(t, result.Domain)
				assert.Equal(t, *tc.expectedDomain, *result.Domain)
			}
		})
	}
}

func TestParse_InvalidInputs(t *testing.T) {
	cases := []string{
		"",
		"scio://",
		"scio://tags/",
		"scio://tags/UPPERCASE",
		"scio://tags/has space",
		"scio://tags/has_underscore",
		"http://tags/wrong-scheme",
		"scio://contexts/",
		"scio://contexts/ecommerce/",
		"scio://contexts/ecommerce/unknown/section",
		"scio://contexts/ecommerce/domains/",
		"scio://contexts/ecommerce/domains/rules/concepts/",
		"scio://contexts/ecommerce/domains/rules/concepts/x/extra",
		"tags/no-scheme",
	}

	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			u, err := uri.Parse(input)
			assert.Error(t, err, "expected error for %q", input)
			assert.Nil(t, u)
		})
	}
}

func TestParentURI(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedParent string
		expectedErr    bool
	}{
		{
			name:        "global tag has no parent",
			input:       "scio://tags/rule",
			expectedErr: true,
		},
		{
			name:        "global relation has no parent",
			input:       "scio://relations/implements",
			expectedErr: true,
		},
		{
			name:        "context has no parent",
			input:       "scio://contexts/ecommerce",
			expectedErr: true,
		},
		{
			name:           "context-scoped tag parent is context",
			input:          "scio://contexts/ecommerce/tags/local-tag",
			expectedParent: "scio://contexts/ecommerce",
		},
		{
			name:           "context-scoped relation parent is context",
			input:          "scio://contexts/ecommerce/relations/owns",
			expectedParent: "scio://contexts/ecommerce",
		},
		{
			name:           "domain parent is context",
			input:          "scio://contexts/ecommerce/domains/business-rules",
			expectedParent: "scio://contexts/ecommerce",
		},
		{
			name:           "concept parent is domain",
			input:          "scio://contexts/ecommerce/domains/business-rules/concepts/discount-calculation",
			expectedParent: "scio://contexts/ecommerce/domains/business-rules",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, err := uri.Parse(tc.input)
			assert.NoError(t, err)

			parent, err := u.ParentURI()

			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedParent, parent)
			}
		})
	}
}

func TestString(t *testing.T) {
	// given
	raw := "scio://tags/critical"
	u, err := uri.Parse(raw)
	assert.NoError(t, err)
	// when
	result := u.String()
	// then
	assert.Equal(t, raw, result)
}
