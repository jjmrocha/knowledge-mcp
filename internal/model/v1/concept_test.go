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
// ParseConcept — happy paths
// ---------------------------------------------------------------------------

func TestParseConcept_ValidConcept(t *testing.T) {
	// given
	meta := `entity: concept
schema: 1
uri: scio://contexts/ecommerce/domains/business-rules/concepts/discount-calculation
name: Discount Calculation
version: 3
created: 2026-02-15T10:00:00Z
last-update: 2026-02-19T14:30:00Z
tags:
    - scio://tags/business-rule
    - scio://tags/pricing
relations:
    - type: scio://relations/implements
      target: scio://contexts/ecommerce/domains/policies/concepts/pricing-policy
sources:
    - type: file
      href: src/pricing/discount.go`
	content := entityContent(meta, "")

	// when
	c, err := v1.ParseConcept(content)

	// then
	require.NoError(t, err)
	require.NotNil(t, c)
	assert.Equal(t, "concept", c.Entity)
	assert.Equal(t, 1, c.Schema)
	assert.Equal(t, "scio://contexts/ecommerce/domains/business-rules/concepts/discount-calculation", c.URI)
	assert.Equal(t, "Discount Calculation", c.Name)
	assert.Equal(t, 3, c.Version)
	assert.True(t, c.Created.Equal(time.Date(2026, 2, 15, 10, 0, 0, 0, time.UTC)))
	assert.True(t, c.LastUpdate.Equal(time.Date(2026, 2, 19, 14, 30, 0, 0, time.UTC)))
	assert.Equal(t, []string{"scio://tags/business-rule", "scio://tags/pricing"}, c.Tags)
	require.Len(t, c.Relations, 1)
	assert.Equal(t, "scio://relations/implements", c.Relations[0].Type)
	assert.Equal(t, "scio://contexts/ecommerce/domains/policies/concepts/pricing-policy", c.Relations[0].Target)
	require.Len(t, c.Sources, 1)
	assert.Equal(t, "file", c.Sources[0].Type)
	assert.Equal(t, "src/pricing/discount.go", c.Sources[0].Href)
}

func TestParseConcept_BodyIsExtracted(t *testing.T) {
	// given
	meta := `entity: concept
schema: 1
uri: scio://contexts/ecommerce/domains/business-rules/concepts/discount-calculation
name: Discount Calculation
version: 1
created: 2026-01-01T00:00:00Z
last-update: 2026-01-01T00:00:00Z
tags: []
relations: []
sources: []`
	body := "Business logic for calculating tiered discounts based on order total.\n"
	content := entityContent(meta, body)

	// when
	c, err := v1.ParseConcept(content)

	// then
	require.NoError(t, err)
	assert.Equal(t, body, c.Body)
}

func TestParseConcept_EmptyTagsRelationsSources(t *testing.T) {
	// given — explicit empty lists
	meta := `entity: concept
schema: 1
uri: scio://contexts/ecommerce/domains/business-rules/concepts/discount-calculation
name: Discount Calculation
version: 1
created: 2026-01-01T00:00:00Z
last-update: 2026-01-01T00:00:00Z
tags: []
relations: []
sources: []`
	content := entityContent(meta, "")

	// when
	c, err := v1.ParseConcept(content)

	// then
	require.NoError(t, err)
	assert.Empty(t, c.Tags)
	assert.Empty(t, c.Relations)
	assert.Empty(t, c.Sources)
}

func TestParseConcept_MultipleSources(t *testing.T) {
	// given
	meta := `entity: concept
schema: 1
uri: scio://contexts/ecommerce/domains/business-rules/concepts/discount-calculation
name: Discount Calculation
version: 1
created: 2026-01-01T00:00:00Z
last-update: 2026-01-01T00:00:00Z
tags: []
relations: []
sources:
    - type: file
      href: src/pricing/discount.go
    - type: url
      href: https://wiki.internal/discount`
	content := entityContent(meta, "")

	// when
	c, err := v1.ParseConcept(content)

	// then
	require.NoError(t, err)
	require.Len(t, c.Sources, 2)
	assert.Equal(t, "file", c.Sources[0].Type)
	assert.Equal(t, "src/pricing/discount.go", c.Sources[0].Href)
	assert.Equal(t, "url", c.Sources[1].Type)
	assert.Equal(t, "https://wiki.internal/discount", c.Sources[1].Href)
}

// ---------------------------------------------------------------------------
// ParseConcept — error paths
// ---------------------------------------------------------------------------

func TestParseConcept_MissingFrontmatter(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"no delimiters", "entity: concept\nschema: 1"},
		{"only opening delimiter", "---\nentity: concept\n"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, err := v1.ParseConcept(tc.input)
			assert.Error(t, err)
			assert.Nil(t, c)
			assert.ErrorContains(t, err, "failed to parse concept file")
		})
	}
}

func TestParseConcept_InvalidYAML(t *testing.T) {
	// given
	content := entityContent("entity: [\nbad yaml", "")

	// when
	c, err := v1.ParseConcept(content)

	// then
	assert.Error(t, err)
	assert.Nil(t, c)
	assert.ErrorContains(t, err, "failed to unmarshal concept metadata")
}

func TestParseConcept_WrongEntityType(t *testing.T) {
	tests := []struct {
		name       string
		entityType string
	}{
		{"tag entity", "tag"},
		{"context entity", "context"},
		{"domain entity", "domain"},
		{"relation entity", "relation"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			meta := "entity: " + tc.entityType + "\nschema: 1\nuri: scio://contexts/x/domains/y/concepts/z\nname: Z\nversion: 1\ncreated: 2026-01-01T00:00:00Z\nlast-update: 2026-01-01T00:00:00Z\ntags: []\nrelations: []\nsources: []"
			content := entityContent(meta, "")

			c, err := v1.ParseConcept(content)

			assert.Error(t, err)
			assert.Nil(t, c)
			assert.ErrorContains(t, err, "invalid entity type")
			assert.ErrorContains(t, err, tc.entityType)
		})
	}
}

// ---------------------------------------------------------------------------
// EncodeConcept — nil-slice defaulting
// ---------------------------------------------------------------------------

func TestEncodeConcept_NilTags_DefaultsToEmpty(t *testing.T) {
	// given
	c := &v1.Concept{
		Entity:     "concept",
		Schema:     1,
		URI:        "scio://contexts/x/domains/y/concepts/z",
		Name:       "Z",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:       nil,
		Relations:  []v1.RelationRef{},
		Sources:    []v1.Source{},
	}

	// when
	encoded, err := v1.EncodeConcept(c)

	// then
	require.NoError(t, err)
	assert.Contains(t, encoded, "tags: []")
}

func TestEncodeConcept_NilRelations_DefaultsToEmpty(t *testing.T) {
	// given
	c := &v1.Concept{
		Entity:     "concept",
		Schema:     1,
		URI:        "scio://contexts/x/domains/y/concepts/z",
		Name:       "Z",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:       []string{},
		Relations:  nil,
		Sources:    []v1.Source{},
	}

	// when
	encoded, err := v1.EncodeConcept(c)

	// then
	require.NoError(t, err)
	assert.Contains(t, encoded, "relations: []")
}

func TestEncodeConcept_NilSources_DefaultsToEmpty(t *testing.T) {
	// given
	c := &v1.Concept{
		Entity:     "concept",
		Schema:     1,
		URI:        "scio://contexts/x/domains/y/concepts/z",
		Name:       "Z",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:       []string{},
		Relations:  []v1.RelationRef{},
		Sources:    nil,
	}

	// when
	encoded, err := v1.EncodeConcept(c)

	// then
	require.NoError(t, err)
	assert.Contains(t, encoded, "sources: []")
}

func TestEncodeConcept_DoesNotMutateOriginal(t *testing.T) {
	// given
	c := &v1.Concept{
		Entity:     "concept",
		Schema:     1,
		URI:        "scio://contexts/x/domains/y/concepts/z",
		Name:       "Z",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		// Tags, Relations, Sources all nil
	}

	// when
	_, err := v1.EncodeConcept(c)

	// then
	require.NoError(t, err)
	assert.Nil(t, c.Tags)
	assert.Nil(t, c.Relations)
	assert.Nil(t, c.Sources)
}

func TestEncodeConcept_ExplicitSlicesPreserved(t *testing.T) {
	// given
	c := &v1.Concept{
		Entity:     "concept",
		Schema:     1,
		URI:        "scio://contexts/x/domains/y/concepts/z",
		Name:       "Z",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:       []string{"scio://tags/business-rule"},
		Relations:  []v1.RelationRef{{Type: "scio://relations/implements", Target: "scio://contexts/x/domains/y/concepts/policy"}},
		Sources:    []v1.Source{{Type: "file", Href: "src/discount.go"}},
	}

	// when
	encoded, err := v1.EncodeConcept(c)

	// then
	require.NoError(t, err)
	assert.Contains(t, encoded, "scio://tags/business-rule")
	assert.Contains(t, encoded, "scio://relations/implements")
	assert.Contains(t, encoded, "src/discount.go")
}

func TestEncodeConcept_BodyIncluded(t *testing.T) {
	// given
	c := &v1.Concept{
		Entity:     "concept",
		Schema:     1,
		URI:        "scio://contexts/x/domains/y/concepts/z",
		Name:       "Z",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:       []string{},
		Relations:  []v1.RelationRef{},
		Sources:    []v1.Source{},
		Body:       "Business logic for calculating tiered discounts.\n",
	}

	// when
	encoded, err := v1.EncodeConcept(c)

	// then
	require.NoError(t, err)
	assert.True(t, strings.HasSuffix(encoded, "Business logic for calculating tiered discounts.\n"))
}

func TestEncodeConcept_OutputStartsWithFrontmatterDelimiter(t *testing.T) {
	// given
	c := &v1.Concept{
		Entity:     "concept",
		Schema:     1,
		URI:        "scio://contexts/x/domains/y/concepts/z",
		Name:       "Z",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// when
	encoded, err := v1.EncodeConcept(c)

	// then
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(encoded, "---\n"))
}

// ---------------------------------------------------------------------------
// Round-trip: EncodeConcept → ParseConcept
// ---------------------------------------------------------------------------

func TestEncodeConcept_ParseConcept_RoundTrip(t *testing.T) {
	// given
	original := &v1.Concept{
		Entity:     "concept",
		Schema:     1,
		URI:        "scio://contexts/ecommerce/domains/business-rules/concepts/discount-calculation",
		Name:       "Discount Calculation",
		Version:    3,
		Created:    time.Date(2026, 2, 15, 10, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 2, 19, 14, 30, 0, 0, time.UTC),
		Tags:       []string{"scio://tags/business-rule", "scio://tags/pricing"},
		Relations: []v1.RelationRef{
			{Type: "scio://relations/implements", Target: "scio://contexts/ecommerce/domains/policies/concepts/pricing-policy"},
		},
		Sources: []v1.Source{
			{Type: "file", Href: "src/pricing/discount.go"},
		},
		Body: "Business logic for calculating tiered discounts based on order total.\n",
	}

	// when
	encoded, err := v1.EncodeConcept(original)
	require.NoError(t, err)

	parsed, err := v1.ParseConcept(encoded)

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
	assert.Equal(t, original.Sources, parsed.Sources)
	assert.Equal(t, original.Body, parsed.Body)
}

func TestEncodeConcept_ParseConcept_NilSlicesRoundTrip(t *testing.T) {
	// given
	original := &v1.Concept{
		Entity:     "concept",
		Schema:     1,
		URI:        "scio://contexts/x/domains/y/concepts/z",
		Name:       "Z",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		// Tags, Relations, Sources all nil
	}

	// when
	encoded, err := v1.EncodeConcept(original)
	require.NoError(t, err)

	parsed, err := v1.ParseConcept(encoded)

	// then
	require.NoError(t, err)
	assert.Empty(t, parsed.Tags)
	assert.Empty(t, parsed.Relations)
	assert.Empty(t, parsed.Sources)
}
