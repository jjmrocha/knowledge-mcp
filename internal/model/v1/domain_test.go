package v1_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "github.com/jjmrocha/knowledge-mcp/internal/model/v1"
)

// ---------------------------------------------------------------------------
// ParseDomain — happy paths
// ---------------------------------------------------------------------------

func TestParseDomain_ValidDomain(t *testing.T) {
	// given
	meta := `entity: domain
schema: 1
uri: scio://contexts/ecommerce/domains/business-rules
name: Business Rules
version: 1
created: 2026-02-15T10:00:00Z
last-update: 2026-02-19T14:30:00Z
tags:
    - scio://tags/core
relations: []`
	content := entityContent(meta, "")

	// when
	d, err := v1.ParseDomain(content)

	// then
	require.NoError(t, err)
	require.NotNil(t, d)
	assert.Equal(t, "domain", d.Entity)
	assert.Equal(t, 1, d.Schema)
	assert.Equal(t, "scio://contexts/ecommerce/domains/business-rules", d.URI)
	assert.Equal(t, "Business Rules", d.Name)
	assert.Equal(t, 1, d.Version)
	assert.True(t, d.Created.Equal(time.Date(2026, 2, 15, 10, 0, 0, 0, time.UTC)))
	assert.True(t, d.LastUpdate.Equal(time.Date(2026, 2, 19, 14, 30, 0, 0, time.UTC)))
	assert.Equal(t, []string{"scio://tags/core"}, d.Tags)
	assert.Empty(t, d.Relations)
}

func TestParseDomain_BodyIsExtracted(t *testing.T) {
	// given
	meta := `entity: domain
schema: 1
uri: scio://contexts/ecommerce/domains/business-rules
name: Business Rules
version: 1
created: 2026-01-01T00:00:00Z
last-update: 2026-01-01T00:00:00Z
tags: []
relations: []`
	body := "Core business rules governing e-commerce operations.\n"
	content := entityContent(meta, body)

	// when
	d, err := v1.ParseDomain(content)

	// then
	require.NoError(t, err)
	assert.Equal(t, body, d.Body)
}

func TestParseDomain_WithRelations(t *testing.T) {
	// given
	meta := `entity: domain
schema: 1
uri: scio://contexts/ecommerce/domains/business-rules
name: Business Rules
version: 1
created: 2026-01-01T00:00:00Z
last-update: 2026-01-01T00:00:00Z
tags: []
relations:
    - type: scio://relations/depends-on
      target: scio://contexts/ecommerce/domains/products`
	content := entityContent(meta, "")

	// when
	d, err := v1.ParseDomain(content)

	// then
	require.NoError(t, err)
	require.Len(t, d.Relations, 1)
	assert.Equal(t, "scio://relations/depends-on", d.Relations[0].Type)
	assert.Equal(t, "scio://contexts/ecommerce/domains/products", d.Relations[0].Target)
}

func TestParseDomain_EmptyTagsAndRelations(t *testing.T) {
	// given
	meta := `entity: domain
schema: 1
uri: scio://contexts/ecommerce/domains/business-rules
name: Business Rules
version: 1
created: 2026-01-01T00:00:00Z
last-update: 2026-01-01T00:00:00Z
tags: []
relations: []`
	content := entityContent(meta, "")

	// when
	d, err := v1.ParseDomain(content)

	// then
	require.NoError(t, err)
	assert.Empty(t, d.Tags)
	assert.Empty(t, d.Relations)
}

// ---------------------------------------------------------------------------
// ParseDomain — error paths
// ---------------------------------------------------------------------------

func TestParseDomain_MissingFrontmatter(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"no delimiters", "entity: domain\nschema: 1"},
		{"only opening delimiter", "---\nentity: domain\n"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d, err := v1.ParseDomain(tc.input)
			assert.Error(t, err)
			assert.Nil(t, d)
			assert.ErrorContains(t, err, "failed to parse domain file")
		})
	}
}

func TestParseDomain_InvalidYAML(t *testing.T) {
	// given
	content := entityContent("entity: [\nbad yaml", "")

	// when
	d, err := v1.ParseDomain(content)

	// then
	assert.Error(t, err)
	assert.Nil(t, d)
	assert.ErrorContains(t, err, "failed to unmarshal domain metadata")
}

func TestParseDomain_WrongEntityType(t *testing.T) {
	tests := []struct {
		name       string
		entityType string
	}{
		{"tag entity", "tag"},
		{"context entity", "context"},
		{"concept entity", "concept"},
		{"relation entity", "relation"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			meta := "entity: " + tc.entityType + "\nschema: 1\nuri: scio://contexts/x/domains/y\nname: Y\nversion: 1\ncreated: 2026-01-01T00:00:00Z\nlast-update: 2026-01-01T00:00:00Z\ntags: []\nrelations: []"
			content := entityContent(meta, "")

			d, err := v1.ParseDomain(content)

			assert.Error(t, err)
			assert.Nil(t, d)
			assert.ErrorContains(t, err, "invalid entity type")
			assert.ErrorContains(t, err, tc.entityType)
		})
	}
}

// ---------------------------------------------------------------------------
// EncodeDomain — nil-slice defaulting
// ---------------------------------------------------------------------------

func TestEncodeDomain_NilTags_DefaultsToEmpty(t *testing.T) {
	// given
	d := &v1.Domain{
		Entity:     "domain",
		Schema:     1,
		URI:        "scio://contexts/ecommerce/domains/business-rules",
		Name:       "Business Rules",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:       nil,
		Relations:  []v1.RelationRef{},
	}

	// when
	encoded, err := v1.EncodeDomain(d)

	// then
	require.NoError(t, err)
	assert.Contains(t, encoded, "tags: []")
}

func TestEncodeDomain_NilRelations_DefaultsToEmpty(t *testing.T) {
	// given
	d := &v1.Domain{
		Entity:     "domain",
		Schema:     1,
		URI:        "scio://contexts/ecommerce/domains/business-rules",
		Name:       "Business Rules",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:       []string{},
		Relations:  nil,
	}

	// when
	encoded, err := v1.EncodeDomain(d)

	// then
	require.NoError(t, err)
	assert.Contains(t, encoded, "relations: []")
}

func TestEncodeDomain_DoesNotMutateOriginal(t *testing.T) {
	// given
	d := &v1.Domain{
		Entity:     "domain",
		Schema:     1,
		URI:        "scio://contexts/x/domains/y",
		Name:       "Y",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// when
	_, err := v1.EncodeDomain(d)

	// then
	require.NoError(t, err)
	assert.Nil(t, d.Tags)
	assert.Nil(t, d.Relations)
}

func TestEncodeDomain_ExplicitSlicesPreserved(t *testing.T) {
	// given
	d := &v1.Domain{
		Entity:     "domain",
		Schema:     1,
		URI:        "scio://contexts/ecommerce/domains/business-rules",
		Name:       "Business Rules",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:       []string{"scio://tags/core"},
		Relations: []v1.RelationRef{
			{Type: "scio://relations/depends-on", Target: "scio://contexts/ecommerce/domains/products"},
		},
	}

	// when
	encoded, err := v1.EncodeDomain(d)

	// then
	require.NoError(t, err)
	assert.Contains(t, encoded, "scio://tags/core")
	assert.Contains(t, encoded, "scio://relations/depends-on")
}

func TestEncodeDomain_BodyIncluded(t *testing.T) {
	// given
	d := &v1.Domain{
		Entity:     "domain",
		Schema:     1,
		URI:        "scio://contexts/x/domains/y",
		Name:       "Y",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:       []string{},
		Relations:  []v1.RelationRef{},
		Body:       "Core business rules governing e-commerce operations.\n",
	}

	// when
	encoded, err := v1.EncodeDomain(d)

	// then
	require.NoError(t, err)
	assert.True(t, strings.HasSuffix(encoded, "Core business rules governing e-commerce operations.\n"))
}

func TestEncodeDomain_OutputStartsWithFrontmatterDelimiter(t *testing.T) {
	// given
	d := &v1.Domain{
		Entity:     "domain",
		Schema:     1,
		URI:        "scio://contexts/x/domains/y",
		Name:       "Y",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// when
	encoded, err := v1.EncodeDomain(d)

	// then
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(encoded, "---\n"))
}

// ---------------------------------------------------------------------------
// Round-trip: EncodeDomain → ParseDomain
// ---------------------------------------------------------------------------

func TestEncodeDomain_ParseDomain_RoundTrip(t *testing.T) {
	// given
	original := &v1.Domain{
		Entity:     "domain",
		Schema:     1,
		URI:        "scio://contexts/ecommerce/domains/business-rules",
		Name:       "Business Rules",
		Version:    3,
		Created:    time.Date(2026, 2, 15, 10, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 2, 19, 14, 30, 0, 0, time.UTC),
		Tags:       []string{"scio://tags/core"},
		Relations: []v1.RelationRef{
			{Type: "scio://relations/depends-on", Target: "scio://contexts/ecommerce/domains/products"},
		},
		Body: "Core business rules.\n",
	}

	// when
	encoded, err := v1.EncodeDomain(original)
	require.NoError(t, err)

	parsed, err := v1.ParseDomain(encoded)

	// then
	require.NoError(t, err)
	require.NotNil(t, parsed)
	assert.Equal(t, original.Entity, parsed.Entity)
	assert.Equal(t, original.Schema, parsed.Schema)
	assert.Equal(t, original.URI, parsed.URI)
	assert.Equal(t, original.Name, parsed.Name)
	assert.Equal(t, original.Version, parsed.Version)
	assert.True(t, original.Created.Equal(parsed.Created))
	assert.True(t, original.LastUpdate.Equal(parsed.LastUpdate))
	assert.Equal(t, original.Tags, parsed.Tags)
	assert.Equal(t, original.Relations, parsed.Relations)
	assert.Equal(t, original.Body, parsed.Body)
}

func TestEncodeDomain_ParseDomain_NilSlicesRoundTrip(t *testing.T) {
	// given
	original := &v1.Domain{
		Entity:     "domain",
		Schema:     1,
		URI:        "scio://contexts/x/domains/y",
		Name:       "Y",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// when
	encoded, err := v1.EncodeDomain(original)
	require.NoError(t, err)

	parsed, err := v1.ParseDomain(encoded)

	// then
	require.NoError(t, err)
	assert.Empty(t, parsed.Tags)
	assert.Empty(t, parsed.Relations)
}
