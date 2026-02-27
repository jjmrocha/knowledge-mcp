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
// ParseContext — happy paths
// ---------------------------------------------------------------------------

func TestParseContext_ValidContext(t *testing.T) {
	// given
	meta := `entity: context
schema: 1
uri: scio://contexts/ecommerce
name: E-Commerce Platform
version: 2
created: 2026-02-15T10:00:00Z
last-update: 2026-02-19T14:30:00Z
tags:
    - scio://tags/business-rule
relations: []`
	content := entityContent(meta, "")

	// when
	ctx, err := v1.ParseContext(content)

	// then
	require.NoError(t, err)
	require.NotNil(t, ctx)
	assert.Equal(t, "context", ctx.Entity)
	assert.Equal(t, 1, ctx.Schema)
	assert.Equal(t, "scio://contexts/ecommerce", ctx.URI)
	assert.Equal(t, "E-Commerce Platform", ctx.Name)
	assert.Equal(t, 2, ctx.Version)
	assert.True(t, ctx.Created.Equal(time.Date(2026, 2, 15, 10, 0, 0, 0, time.UTC)))
	assert.True(t, ctx.LastUpdate.Equal(time.Date(2026, 2, 19, 14, 30, 0, 0, time.UTC)))
	assert.Equal(t, []string{"scio://tags/business-rule"}, ctx.Tags)
	assert.Empty(t, ctx.Relations)
}

func TestParseContext_BodyIsExtracted(t *testing.T) {
	// given
	meta := `entity: context
schema: 1
uri: scio://contexts/ecommerce
name: E-Commerce Platform
version: 1
created: 2026-01-01T00:00:00Z
last-update: 2026-01-01T00:00:00Z
tags: []
relations: []`
	body := "Top-level knowledge area for e-commerce domain.\n"
	content := entityContent(meta, body)

	// when
	ctx, err := v1.ParseContext(content)

	// then
	require.NoError(t, err)
	assert.Equal(t, body, ctx.Body)
}

func TestParseContext_WithRelations(t *testing.T) {
	// given
	meta := `entity: context
schema: 1
uri: scio://contexts/ecommerce
name: E-Commerce Platform
version: 1
created: 2026-01-01T00:00:00Z
last-update: 2026-01-01T00:00:00Z
tags: []
relations:
    - type: scio://relations/depends-on
      target: scio://contexts/payment`
	content := entityContent(meta, "")

	// when
	ctx, err := v1.ParseContext(content)

	// then
	require.NoError(t, err)
	require.Len(t, ctx.Relations, 1)
	assert.Equal(t, "scio://relations/depends-on", ctx.Relations[0].Type)
	assert.Equal(t, "scio://contexts/payment", ctx.Relations[0].Target)
}

func TestParseContext_EmptyTagsAndRelations(t *testing.T) {
	// given — explicit empty lists
	meta := `entity: context
schema: 1
uri: scio://contexts/ecommerce
name: E-Commerce
version: 1
created: 2026-01-01T00:00:00Z
last-update: 2026-01-01T00:00:00Z
tags: []
relations: []`
	content := entityContent(meta, "")

	// when
	ctx, err := v1.ParseContext(content)

	// then
	require.NoError(t, err)
	assert.Empty(t, ctx.Tags)
	assert.Empty(t, ctx.Relations)
}

// ---------------------------------------------------------------------------
// ParseContext — error paths
// ---------------------------------------------------------------------------

func TestParseContext_MissingFrontmatter(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"no delimiters", "entity: context\nschema: 1"},
		{"only opening delimiter", "---\nentity: context\n"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, err := v1.ParseContext(tc.input)
			assert.Error(t, err)
			assert.Nil(t, ctx)
			assert.ErrorContains(t, err, "failed to parse context file")
		})
	}
}

func TestParseContext_InvalidYAML(t *testing.T) {
	// given
	content := entityContent("entity: [\nbad yaml", "")

	// when
	ctx, err := v1.ParseContext(content)

	// then
	assert.Error(t, err)
	assert.Nil(t, ctx)
	assert.ErrorContains(t, err, "failed to unmarshal context metadata")
}

func TestParseContext_WrongEntityType(t *testing.T) {
	tests := []struct {
		name       string
		entityType string
	}{
		{"tag entity", "tag"},
		{"domain entity", "domain"},
		{"concept entity", "concept"},
		{"relation entity", "relation"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			meta := "entity: " + tc.entityType + "\nschema: 1\nuri: scio://contexts/x\nname: X\nversion: 1\ncreated: 2026-01-01T00:00:00Z\nlast-update: 2026-01-01T00:00:00Z\ntags: []\nrelations: []"
			content := entityContent(meta, "")

			ctx, err := v1.ParseContext(content)

			assert.Error(t, err)
			assert.Nil(t, ctx)
			assert.ErrorContains(t, err, "invalid entity type")
			assert.ErrorContains(t, err, tc.entityType)
		})
	}
}

// ---------------------------------------------------------------------------
// EncodeContext — nil-slice defaulting
// ---------------------------------------------------------------------------

func TestEncodeContext_NilTags_DefaultsToEmpty(t *testing.T) {
	// given
	ctx := &v1.Context{
		Entity:     "context",
		Schema:     1,
		URI:        "scio://contexts/ecommerce",
		Name:       "E-Commerce",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:       nil,
		Relations:  []v1.RelationRef{},
	}

	// when
	encoded, err := v1.EncodeContext(ctx)

	// then
	require.NoError(t, err)
	assert.Contains(t, encoded, "tags: []")
}

func TestEncodeContext_NilRelations_DefaultsToEmpty(t *testing.T) {
	// given
	ctx := &v1.Context{
		Entity:     "context",
		Schema:     1,
		URI:        "scio://contexts/ecommerce",
		Name:       "E-Commerce",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:       []string{},
		Relations:  nil,
	}

	// when
	encoded, err := v1.EncodeContext(ctx)

	// then
	require.NoError(t, err)
	assert.Contains(t, encoded, "relations: []")
}

func TestEncodeContext_DoesNotMutateOriginal(t *testing.T) {
	// given
	ctx := &v1.Context{
		Entity:     "context",
		Schema:     1,
		URI:        "scio://contexts/x",
		Name:       "X",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		// Tags and Relations nil — defaults are applied to copy only
	}

	// when
	_, err := v1.EncodeContext(ctx)

	// then
	require.NoError(t, err)
	assert.Nil(t, ctx.Tags)
	assert.Nil(t, ctx.Relations)
}

func TestEncodeContext_ExplicitSlicesPreserved(t *testing.T) {
	// given
	ctx := &v1.Context{
		Entity:     "context",
		Schema:     1,
		URI:        "scio://contexts/ecommerce",
		Name:       "E-Commerce",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:       []string{"scio://tags/business-rule"},
		Relations: []v1.RelationRef{
			{Type: "scio://relations/depends-on", Target: "scio://contexts/payment"},
		},
	}

	// when
	encoded, err := v1.EncodeContext(ctx)

	// then
	require.NoError(t, err)
	assert.Contains(t, encoded, "scio://tags/business-rule")
	assert.Contains(t, encoded, "scio://relations/depends-on")
	assert.Contains(t, encoded, "scio://contexts/payment")
}

func TestEncodeContext_BodyIncluded(t *testing.T) {
	// given
	ctx := &v1.Context{
		Entity:     "context",
		Schema:     1,
		URI:        "scio://contexts/x",
		Name:       "X",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:       []string{},
		Relations:  []v1.RelationRef{},
		Body:       "Top-level knowledge area.\n",
	}

	// when
	encoded, err := v1.EncodeContext(ctx)

	// then
	require.NoError(t, err)
	assert.True(t, strings.HasSuffix(encoded, "Top-level knowledge area.\n"))
}

func TestEncodeContext_OutputStartsWithFrontmatterDelimiter(t *testing.T) {
	// given
	ctx := &v1.Context{
		Entity:     "context",
		Schema:     1,
		URI:        "scio://contexts/x",
		Name:       "X",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// when
	encoded, err := v1.EncodeContext(ctx)

	// then
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(encoded, "---\n"))
}

// ---------------------------------------------------------------------------
// Round-trip: EncodeContext → ParseContext
// ---------------------------------------------------------------------------

func TestEncodeContext_ParseContext_RoundTrip(t *testing.T) {
	// given
	original := &v1.Context{
		Entity:     "context",
		Schema:     1,
		URI:        "scio://contexts/ecommerce",
		Name:       "E-Commerce Platform",
		Version:    2,
		Created:    time.Date(2026, 2, 15, 10, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 2, 19, 14, 30, 0, 0, time.UTC),
		Tags:       []string{"scio://tags/business-rule"},
		Relations: []v1.RelationRef{
			{Type: "scio://relations/depends-on", Target: "scio://contexts/payment"},
		},
		Body: "Top-level e-commerce knowledge area.\n",
	}

	// when
	encoded, err := v1.EncodeContext(original)
	require.NoError(t, err)

	parsed, err := v1.ParseContext(encoded)

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

func TestEncodeContext_ParseContext_NilSlicesRoundTrip(t *testing.T) {
	// given — nil slices trigger defaulting; verify the encoded defaults round-trip cleanly
	original := &v1.Context{
		Entity:     "context",
		Schema:     1,
		URI:        "scio://contexts/x",
		Name:       "X",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		// Tags and Relations nil
	}

	// when
	encoded, err := v1.EncodeContext(original)
	require.NoError(t, err)

	parsed, err := v1.ParseContext(encoded)

	// then
	require.NoError(t, err)
	assert.Empty(t, parsed.Tags)
	assert.Empty(t, parsed.Relations)
}
