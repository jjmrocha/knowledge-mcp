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
// ParseRelationType — happy paths
// ---------------------------------------------------------------------------

func TestParseRelationType_ValidRelation(t *testing.T) {
	// given
	meta := `entity: relation
schema: 1
uri: scio://relations/implements
version: 1
created: 2026-02-15T10:00:00Z
last-update: 2026-02-19T14:30:00Z
inverse-of: scio://relations/implemented-by
allowed-source-entities:
    - context
    - domain
    - concept
allowed-target-entities:
    - concept
transitive: false
symmetric: false`
	content := entityContent(meta, "")

	// when
	r, err := v1.ParseRelationType(content)

	// then
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "relation", r.Entity)
	assert.Equal(t, 1, r.Schema)
	assert.Equal(t, "scio://relations/implements", r.URI)
	assert.Equal(t, 1, r.Version)
	assert.True(t, r.Created.Equal(time.Date(2026, 2, 15, 10, 0, 0, 0, time.UTC)))
	assert.True(t, r.LastUpdate.Equal(time.Date(2026, 2, 19, 14, 30, 0, 0, time.UTC)))
	assert.Equal(t, "scio://relations/implemented-by", r.InverseOf)
	assert.Equal(t, []string{"context", "domain", "concept"}, r.AllowedSourceEntities)
	assert.Equal(t, []string{"concept"}, r.AllowedTargetEntities)
	assert.False(t, r.Transitive)
	assert.False(t, r.Symmetric)
}

func TestParseRelationType_TransitiveAndSymmetric(t *testing.T) {
	// given
	meta := `entity: relation
schema: 1
uri: scio://relations/related-to
version: 1
created: 2026-01-01T00:00:00Z
last-update: 2026-01-01T00:00:00Z
allowed-source-entities: []
allowed-target-entities: []
transitive: true
symmetric: true`
	content := entityContent(meta, "")

	// when
	r, err := v1.ParseRelationType(content)

	// then
	require.NoError(t, err)
	assert.True(t, r.Transitive)
	assert.True(t, r.Symmetric)
}

func TestParseRelationType_NoInverseOf(t *testing.T) {
	// given — omitempty: missing inverse-of is valid
	meta := `entity: relation
schema: 1
uri: scio://relations/owns
version: 1
created: 2026-01-01T00:00:00Z
last-update: 2026-01-01T00:00:00Z
allowed-source-entities: []
allowed-target-entities: []
transitive: false
symmetric: false`
	content := entityContent(meta, "")

	// when
	r, err := v1.ParseRelationType(content)

	// then
	require.NoError(t, err)
	assert.Empty(t, r.InverseOf)
}

func TestParseRelationType_BodyIsExtracted(t *testing.T) {
	// given
	meta := `entity: relation
schema: 1
uri: scio://relations/implements
version: 1
created: 2026-01-01T00:00:00Z
last-update: 2026-01-01T00:00:00Z
allowed-source-entities: []
allowed-target-entities: []
transitive: false
symmetric: false`
	body := "Indicates that an entity implements a policy or specification.\n"
	content := entityContent(meta, body)

	// when
	r, err := v1.ParseRelationType(content)

	// then
	require.NoError(t, err)
	assert.Equal(t, body, r.Body)
}

// ---------------------------------------------------------------------------
// ParseRelationType — error paths
// ---------------------------------------------------------------------------

func TestParseRelationType_MissingFrontmatter(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"no delimiters", "entity: relation\nschema: 1"},
		{"only opening delimiter", "---\nentity: relation\n"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, err := v1.ParseRelationType(tc.input)
			assert.Error(t, err)
			assert.Nil(t, r)
			assert.ErrorContains(t, err, "failed to parse relation file")
		})
	}
}

func TestParseRelationType_InvalidYAML(t *testing.T) {
	// given
	content := entityContent("entity: [\nbad yaml", "")

	// when
	r, err := v1.ParseRelationType(content)

	// then
	assert.Error(t, err)
	assert.Nil(t, r)
	assert.ErrorContains(t, err, "failed to unmarshal relation metadata")
}

func TestParseRelationType_WrongEntityType(t *testing.T) {
	tests := []struct {
		name       string
		entityType string
	}{
		{"tag entity", "tag"},
		{"context entity", "context"},
		{"domain entity", "domain"},
		{"concept entity", "concept"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			meta := "entity: " + tc.entityType + "\nschema: 1\nuri: scio://relations/x\nversion: 1\ncreated: 2026-01-01T00:00:00Z\nlast-update: 2026-01-01T00:00:00Z\nallowed-source-entities: []\nallowed-target-entities: []\ntransitive: false\nsymmetric: false"
			content := entityContent(meta, "")

			r, err := v1.ParseRelationType(content)

			assert.Error(t, err)
			assert.Nil(t, r)
			assert.ErrorContains(t, err, "invalid entity type")
			assert.ErrorContains(t, err, tc.entityType)
		})
	}
}

// ---------------------------------------------------------------------------
// EncodeRelationType — nil-slice defaulting
// ---------------------------------------------------------------------------

func TestEncodeRelationType_NilAllowedSourceEntities_DefaultsToAll(t *testing.T) {
	// given
	r := &v1.RelationType{
		Entity:                "relation",
		Schema:                1,
		URI:                   "scio://relations/implements",
		Version:               1,
		Created:               time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate:            time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		AllowedSourceEntities: nil,
		AllowedTargetEntities: []string{},
	}

	// when
	encoded, err := v1.EncodeRelationType(r)

	// then
	require.NoError(t, err)
	assert.Contains(t, encoded, "- context")
	assert.Contains(t, encoded, "- domain")
	assert.Contains(t, encoded, "- concept")
}

func TestEncodeRelationType_NilAllowedTargetEntities_DefaultsToAll(t *testing.T) {
	// given
	r := &v1.RelationType{
		Entity:                "relation",
		Schema:                1,
		URI:                   "scio://relations/implements",
		Version:               1,
		Created:               time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate:            time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		AllowedSourceEntities: []string{},
		AllowedTargetEntities: nil,
	}

	// when
	encoded, err := v1.EncodeRelationType(r)

	// then
	require.NoError(t, err)
	assert.Contains(t, encoded, "- context")
	assert.Contains(t, encoded, "- domain")
	assert.Contains(t, encoded, "- concept")
}

func TestEncodeRelationType_DoesNotMutateOriginal(t *testing.T) {
	// given
	r := &v1.RelationType{
		Entity:     "relation",
		Schema:     1,
		URI:        "scio://relations/x",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		// AllowedSourceEntities and AllowedTargetEntities nil
	}

	// when
	_, err := v1.EncodeRelationType(r)

	// then
	require.NoError(t, err)
	assert.Nil(t, r.AllowedSourceEntities)
	assert.Nil(t, r.AllowedTargetEntities)
}

func TestEncodeRelationType_ExplicitSlicesPreserved(t *testing.T) {
	// given — restricted to concept only
	r := &v1.RelationType{
		Entity:                "relation",
		Schema:                1,
		URI:                   "scio://relations/implements",
		Version:               1,
		Created:               time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate:            time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		AllowedSourceEntities: []string{"concept"},
		AllowedTargetEntities: []string{"concept"},
	}

	// when
	encoded, err := v1.EncodeRelationType(r)

	// then
	require.NoError(t, err)
	assert.Contains(t, encoded, "- concept")
	assert.NotContains(t, encoded, "- context")
	assert.NotContains(t, encoded, "- domain")
}

func TestEncodeRelationType_InverseOfOmittedWhenEmpty(t *testing.T) {
	// given
	r := &v1.RelationType{
		Entity:                "relation",
		Schema:                1,
		URI:                   "scio://relations/owns",
		Version:               1,
		Created:               time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate:            time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		InverseOf:             "",
		AllowedSourceEntities: []string{},
		AllowedTargetEntities: []string{},
	}

	// when
	encoded, err := v1.EncodeRelationType(r)

	// then
	require.NoError(t, err)
	assert.NotContains(t, encoded, "inverse-of")
}

func TestEncodeRelationType_InverseOfIncludedWhenSet(t *testing.T) {
	// given
	r := &v1.RelationType{
		Entity:                "relation",
		Schema:                1,
		URI:                   "scio://relations/implements",
		Version:               1,
		Created:               time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate:            time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		InverseOf:             "scio://relations/implemented-by",
		AllowedSourceEntities: []string{},
		AllowedTargetEntities: []string{},
	}

	// when
	encoded, err := v1.EncodeRelationType(r)

	// then
	require.NoError(t, err)
	assert.Contains(t, encoded, "scio://relations/implemented-by")
}

func TestEncodeRelationType_BodyIncluded(t *testing.T) {
	// given
	r := &v1.RelationType{
		Entity:                "relation",
		Schema:                1,
		URI:                   "scio://relations/implements",
		Version:               1,
		Created:               time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate:            time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		AllowedSourceEntities: []string{},
		AllowedTargetEntities: []string{},
		Body:                  "Indicates that an entity implements a specification.\n",
	}

	// when
	encoded, err := v1.EncodeRelationType(r)

	// then
	require.NoError(t, err)
	assert.True(t, strings.HasSuffix(encoded, "Indicates that an entity implements a specification.\n"))
}

func TestEncodeRelationType_OutputStartsWithFrontmatterDelimiter(t *testing.T) {
	// given
	r := &v1.RelationType{
		Entity:     "relation",
		Schema:     1,
		URI:        "scio://relations/x",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// when
	encoded, err := v1.EncodeRelationType(r)

	// then
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(encoded, "---\n"))
}

// ---------------------------------------------------------------------------
// Round-trip: EncodeRelationType → ParseRelationType
// ---------------------------------------------------------------------------

func TestEncodeRelationType_ParseRelationType_RoundTrip(t *testing.T) {
	// given
	original := &v1.RelationType{
		Entity:                "relation",
		Schema:                1,
		URI:                   "scio://relations/implements",
		Version:               2,
		Created:               time.Date(2026, 2, 15, 10, 0, 0, 0, time.UTC),
		LastUpdate:            time.Date(2026, 2, 19, 14, 30, 0, 0, time.UTC),
		InverseOf:             "scio://relations/implemented-by",
		AllowedSourceEntities: []string{"context", "domain", "concept"},
		AllowedTargetEntities: []string{"concept"},
		Transitive:            false,
		Symmetric:             false,
		Body:                  "Indicates that an entity implements a specification.\n",
	}

	// when
	encoded, err := v1.EncodeRelationType(original)
	require.NoError(t, err)

	parsed, err := v1.ParseRelationType(encoded)

	// then
	require.NoError(t, err)
	require.NotNil(t, parsed)
	assert.Equal(t, original.Entity, parsed.Entity)
	assert.Equal(t, original.Schema, parsed.Schema)
	assert.Equal(t, original.URI, parsed.URI)
	assert.Equal(t, original.Version, parsed.Version)
	assert.True(t, original.Created.Equal(parsed.Created))
	assert.True(t, original.LastUpdate.Equal(parsed.LastUpdate))
	assert.Equal(t, original.InverseOf, parsed.InverseOf)
	assert.Equal(t, original.AllowedSourceEntities, parsed.AllowedSourceEntities)
	assert.Equal(t, original.AllowedTargetEntities, parsed.AllowedTargetEntities)
	assert.Equal(t, original.Transitive, parsed.Transitive)
	assert.Equal(t, original.Symmetric, parsed.Symmetric)
	assert.Equal(t, original.Body, parsed.Body)
}

func TestEncodeRelationType_ParseRelationType_NilSlicesRoundTrip(t *testing.T) {
	// given — nil slices trigger defaulting to all entity types
	original := &v1.RelationType{
		Entity:     "relation",
		Schema:     1,
		URI:        "scio://relations/x",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		// AllowedSourceEntities and AllowedTargetEntities nil → default to [context, domain, concept]
	}

	// when
	encoded, err := v1.EncodeRelationType(original)
	require.NoError(t, err)

	parsed, err := v1.ParseRelationType(encoded)

	// then
	require.NoError(t, err)
	assert.Equal(t, []string{"context", "domain", "concept"}, parsed.AllowedSourceEntities)
	assert.Equal(t, []string{"context", "domain", "concept"}, parsed.AllowedTargetEntities)
}
