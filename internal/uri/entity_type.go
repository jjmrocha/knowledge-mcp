package uri

import "regexp"

type EntityType string

const (
	EntityTypeTag      EntityType = "tag"
	EntityTypeRelation EntityType = "relation"
	EntityTypeContext  EntityType = "context"
	EntityTypeDomain   EntityType = "domain"
	EntityTypeConcept  EntityType = "concept"
)

type uriPattern struct {
	re         *regexp.Regexp
	entityType EntityType
	hasContext bool
	hasDomain  bool
}

var uriPatterns = []uriPattern{
	{
		re:         regexp.MustCompile(`^scio://tags/([a-z0-9-]+)$`),
		entityType: EntityTypeTag,
	},
	{
		re:         regexp.MustCompile(`^scio://contexts/([a-z0-9-]+)/tags/([a-z0-9-]+)$`),
		entityType: EntityTypeTag,
		hasContext: true,
	},
	{
		re:         regexp.MustCompile(`^scio://relations/([a-z0-9-]+)$`),
		entityType: EntityTypeRelation,
	},
	{
		re:         regexp.MustCompile(`^scio://contexts/([a-z0-9-]+)/relations/([a-z0-9-]+)$`),
		entityType: EntityTypeRelation,
		hasContext: true,
	},
	{
		re:         regexp.MustCompile(`^scio://contexts/([a-z0-9-]+)/domains/([a-z0-9-]+)/concepts/([a-z0-9-]+)$`),
		entityType: EntityTypeConcept,
		hasContext: true,
		hasDomain:  true,
	},
	{
		re:         regexp.MustCompile(`^scio://contexts/([a-z0-9-]+)/domains/([a-z0-9-]+)$`),
		entityType: EntityTypeDomain,
		hasContext: true,
	},
	{
		re:         regexp.MustCompile(`^scio://contexts/([a-z0-9-]+)$`),
		entityType: EntityTypeContext,
	},
}
