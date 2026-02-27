package model_test

import (
	"strings"
	"testing"
	"time"

	v1 "github.com/jjmrocha/knowledge-mcp/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// tagContent builds a minimal valid tag file content string.
func tagContent(meta, body string) string {
	return "---\n" + meta + "\n---\n" + body
}

// ---------------------------------------------------------------------------
// ParseTag — happy paths
// ---------------------------------------------------------------------------

func TestParseTag_ValidTag(t *testing.T) {
	// given
	meta := `entity: tag
schema: 1
uri: scio://tags/business-rule
version: 1
created: 2026-02-15T10:00:00Z
last-update: 2026-02-19T14:30:00Z
allowed-entities:
    - context
    - domain
    - concept
broader: []
narrower: []`
	content := tagContent(meta, "")

	// when
	tag, err := v1.ParseTag(content)

	// then
	require.NoError(t, err)
	require.NotNil(t, tag)
	assert.Equal(t, "tag", tag.Entity)
	assert.Equal(t, 1, tag.Schema)
	assert.Equal(t, "scio://tags/business-rule", tag.URI)
	assert.Equal(t, 1, tag.Version)
	assert.True(t, tag.Created.Equal(time.Date(2026, 2, 15, 10, 0, 0, 0, time.UTC)))
	assert.True(t, tag.LastUpdate.Equal(time.Date(2026, 2, 19, 14, 30, 0, 0, time.UTC)))
	assert.Equal(t, []string{"context", "domain", "concept"}, tag.AllowedEntities)
	assert.Equal(t, []string{}, tag.Broader)
	assert.Equal(t, []string{}, tag.Narrower)
}

func TestParseTag_BodyIsExtracted(t *testing.T) {
	// given
	meta := `entity: tag
schema: 1
uri: scio://tags/pricing
version: 2
created: 2026-01-01T00:00:00Z
last-update: 2026-01-02T00:00:00Z
allowed-entities:
    - concept
broader: []
narrower: []`
	body := "Tags all entities related to pricing logic.\n"
	content := tagContent(meta, body)

	// when
	tag, err := v1.ParseTag(content)

	// then
	require.NoError(t, err)
	assert.Equal(t, body, tag.Body)
}

func TestParseTag_EmptyAllowedEntities(t *testing.T) {
	// given — explicit empty list disables the tag
	meta := `entity: tag
schema: 1
uri: scio://tags/disabled
version: 1
created: 2026-01-01T00:00:00Z
last-update: 2026-01-01T00:00:00Z
allowed-entities: []
broader: []
narrower: []`
	content := tagContent(meta, "")

	// when
	tag, err := v1.ParseTag(content)

	// then
	require.NoError(t, err)
	assert.Empty(t, tag.AllowedEntities)
}

func TestParseTag_BroaderAndNarrower(t *testing.T) {
	// given
	meta := `entity: tag
schema: 1
uri: scio://tags/discount
version: 1
created: 2026-01-01T00:00:00Z
last-update: 2026-01-01T00:00:00Z
allowed-entities:
    - concept
broader:
    - scio://tags/pricing
narrower:
    - scio://tags/seasonal-discount
    - scio://tags/bulk-discount`
	content := tagContent(meta, "")

	// when
	tag, err := v1.ParseTag(content)

	// then
	require.NoError(t, err)
	assert.Equal(t, []string{"scio://tags/pricing"}, tag.Broader)
	assert.Equal(t, []string{"scio://tags/seasonal-discount", "scio://tags/bulk-discount"}, tag.Narrower)
}

// ---------------------------------------------------------------------------
// ParseTag — error paths
// ---------------------------------------------------------------------------

func TestParseTag_MissingFrontmatter(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"no delimiters", "entity: tag\nschema: 1"},
		{"only opening delimiter", "---\nentity: tag\n"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tag, err := v1.ParseTag(tc.input)
			assert.Error(t, err)
			assert.Nil(t, tag)
			assert.ErrorContains(t, err, "failed to parse tag file")
		})
	}
}

func TestParseTag_InvalidYAML(t *testing.T) {
	// given — structurally broken YAML
	content := tagContent("entity: [\nbad yaml", "")

	// when
	tag, err := v1.ParseTag(content)

	// then
	assert.Error(t, err)
	assert.Nil(t, tag)
	assert.ErrorContains(t, err, "failed to unmarshal tag metadata")
}

func TestParseTag_WrongEntityType(t *testing.T) {
	tests := []struct {
		name       string
		entityType string
	}{
		{"concept entity", "concept"},
		{"context entity", "context"},
		{"domain entity", "domain"},
		{"relation entity", "relation"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			meta := "entity: " + tc.entityType + "\nschema: 1\nuri: scio://tags/x\nversion: 1\ncreated: 2026-01-01T00:00:00Z\nlast-update: 2026-01-01T00:00:00Z\nallowed-entities: []\nbroader: []\nnarrower: []"
			content := tagContent(meta, "")

			tag, err := v1.ParseTag(content)

			assert.Error(t, err)
			assert.Nil(t, tag)
			assert.ErrorContains(t, err, "invalid entity type")
			assert.ErrorContains(t, err, tc.entityType)
		})
	}
}

// ---------------------------------------------------------------------------
// EncodeTag — nil-slice defaulting
// ---------------------------------------------------------------------------

func TestEncodeTag_NilAllowedEntities_DefaultsToAll(t *testing.T) {
	// given — AllowedEntities is nil
	tag := &v1.Tag{
		Entity:          "tag",
		Schema:          1,
		URI:             "scio://tags/new",
		Version:         1,
		Created:         time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		AllowedEntities: nil,
		Broader:         []string{},
		Narrower:        []string{},
	}

	// when
	encoded, err := v1.EncodeTag(tag)

	// then
	require.NoError(t, err)
	assert.Contains(t, encoded, "- context")
	assert.Contains(t, encoded, "- domain")
	assert.Contains(t, encoded, "- concept")
}

func TestEncodeTag_NilAllowedEntities_DoesNotMutateOriginal(t *testing.T) {
	// given
	tag := &v1.Tag{
		Entity:     "tag",
		Schema:     1,
		URI:        "scio://tags/x",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// when
	_, err := v1.EncodeTag(tag)

	// then — original struct must not be modified
	require.NoError(t, err)
	assert.Nil(t, tag.AllowedEntities)
	assert.Nil(t, tag.Broader)
	assert.Nil(t, tag.Narrower)
}

func TestEncodeTag_NilBroader_DefaultsToEmpty(t *testing.T) {
	// given
	tag := &v1.Tag{
		Entity:          "tag",
		Schema:          1,
		URI:             "scio://tags/x",
		Version:         1,
		Created:         time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		AllowedEntities: []string{},
		Broader:         nil,
		Narrower:        []string{},
	}

	// when
	encoded, err := v1.EncodeTag(tag)

	// then — "broader: []" must appear in the encoded output
	require.NoError(t, err)
	assert.Contains(t, encoded, "broader: []")
}

func TestEncodeTag_NilNarrower_DefaultsToEmpty(t *testing.T) {
	// given
	tag := &v1.Tag{
		Entity:          "tag",
		Schema:          1,
		URI:             "scio://tags/x",
		Version:         1,
		Created:         time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		AllowedEntities: []string{},
		Broader:         []string{},
		Narrower:        nil,
	}

	// when
	encoded, err := v1.EncodeTag(tag)

	// then — "narrower: []" must appear in the encoded output
	require.NoError(t, err)
	assert.Contains(t, encoded, "narrower: []")
}

func TestEncodeTag_ExplicitSlicesPreserved(t *testing.T) {
	// given — explicit non-nil slices must not be replaced by defaults
	tag := &v1.Tag{
		Entity:          "tag",
		Schema:          1,
		URI:             "scio://tags/discount",
		Version:         1,
		Created:         time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		AllowedEntities: []string{"concept"},
		Broader:         []string{"scio://tags/pricing"},
		Narrower:        []string{"scio://tags/seasonal-discount"},
	}

	// when
	encoded, err := v1.EncodeTag(tag)

	// then
	require.NoError(t, err)
	assert.Contains(t, encoded, "- concept")
	assert.NotContains(t, encoded, "- context")
	assert.Contains(t, encoded, "scio://tags/pricing")
	assert.Contains(t, encoded, "scio://tags/seasonal-discount")
}

func TestEncodeTag_BodyIncluded(t *testing.T) {
	// given
	tag := &v1.Tag{
		Entity:          "tag",
		Schema:          1,
		URI:             "scio://tags/x",
		Version:         1,
		Created:         time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		AllowedEntities: []string{},
		Broader:         []string{},
		Narrower:        []string{},
		Body:            "Human-readable description of the tag.\n",
	}

	// when
	encoded, err := v1.EncodeTag(tag)

	// then
	require.NoError(t, err)
	assert.True(t, strings.HasSuffix(encoded, "Human-readable description of the tag.\n"))
}

func TestEncodeTag_OutputStartsWithFrontmatterDelimiter(t *testing.T) {
	// given
	tag := &v1.Tag{
		Entity:     "tag",
		Schema:     1,
		URI:        "scio://tags/x",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// when
	encoded, err := v1.EncodeTag(tag)

	// then
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(encoded, "---\n"))
}

// ---------------------------------------------------------------------------
// Round-trip: EncodeTag → ParseTag
// ---------------------------------------------------------------------------

func TestEncodeTag_ParseTag_RoundTrip(t *testing.T) {
	// given
	original := &v1.Tag{
		Entity:          "tag",
		Schema:          1,
		URI:             "scio://tags/business-rule",
		Version:         3,
		Created:         time.Date(2026, 2, 15, 10, 0, 0, 0, time.UTC),
		LastUpdate:      time.Date(2026, 2, 19, 14, 30, 0, 0, time.UTC),
		AllowedEntities: []string{"context", "domain", "concept"},
		Broader:         []string{"scio://tags/rule"},
		Narrower:        []string{"scio://tags/payment-rule"},
		Body:            "All entities that encode a business rule.\n",
	}

	// when
	encoded, err := v1.EncodeTag(original)
	require.NoError(t, err)

	parsed, err := v1.ParseTag(encoded)

	// then
	require.NoError(t, err)
	require.NotNil(t, parsed)
	assert.Equal(t, original.Entity, parsed.Entity)
	assert.Equal(t, original.Schema, parsed.Schema)
	assert.Equal(t, original.URI, parsed.URI)
	assert.Equal(t, original.Version, parsed.Version)
	assert.True(t, original.Created.Equal(parsed.Created))
	assert.True(t, original.LastUpdate.Equal(parsed.LastUpdate))
	assert.Equal(t, original.AllowedEntities, parsed.AllowedEntities)
	assert.Equal(t, original.Broader, parsed.Broader)
	assert.Equal(t, original.Narrower, parsed.Narrower)
	assert.Equal(t, original.Body, parsed.Body)
}

func TestEncodeTag_ParseTag_NilSlicesRoundTrip(t *testing.T) {
	// given — nil slices trigger defaulting; verify the encoded defaults round-trip cleanly
	original := &v1.Tag{
		Entity:     "tag",
		Schema:     1,
		URI:        "scio://tags/x",
		Version:    1,
		Created:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		LastUpdate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		// AllowedEntities, Broader, Narrower all nil → defaults applied during encode
	}

	// when
	encoded, err := v1.EncodeTag(original)
	require.NoError(t, err)

	parsed, err := v1.ParseTag(encoded)

	// then — defaults applied during encode must survive the parse
	require.NoError(t, err)
	assert.Equal(t, []string{"context", "domain", "concept"}, parsed.AllowedEntities)
	assert.Equal(t, []string{}, parsed.Broader)
	assert.Equal(t, []string{}, parsed.Narrower)
}
