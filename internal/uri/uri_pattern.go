package uri

import (
	"regexp"

	"github.com/jjmrocha/knowledge-mcp/internal/model"
)

type uriPattern struct {
	re         *regexp.Regexp
	entityType string
	hasContext bool
	hasDomain  bool
}

var uriPatterns = []uriPattern{
	{
		re:         regexp.MustCompile(`^scio://tags/([a-z]+[a-z0-9-]*)$`),
		entityType: model.EntityTypeTag,
	},
	{
		re:         regexp.MustCompile(`^scio://contexts/([a-z]+[a-z0-9-]*)/tags/([a-z]+[a-z0-9-]*)$`),
		entityType: model.EntityTypeTag,
		hasContext: true,
	},
	{
		re:         regexp.MustCompile(`^scio://relations/([a-z]+[a-z0-9-]*)$`),
		entityType: model.EntityTypeRelation,
	},
	{
		re:         regexp.MustCompile(`^scio://contexts/([a-z]+[a-z0-9-]*)/relations/([a-z]+[a-z0-9-]*)$`),
		entityType: model.EntityTypeRelation,
		hasContext: true,
	},
	{
		re:         regexp.MustCompile(`^scio://contexts/([a-z]+[a-z0-9-]*)/domains/([a-z]+[a-z0-9-]*)/concepts/([a-z]+[a-z0-9-]*)$`),
		entityType: model.EntityTypeConcept,
		hasContext: true,
		hasDomain:  true,
	},
	{
		re:         regexp.MustCompile(`^scio://contexts/([a-z]+[a-z0-9-]*)/domains/([a-z]+[a-z0-9-]*)$`),
		entityType: model.EntityTypeDomain,
		hasContext: true,
	},
	{
		re:         regexp.MustCompile(`^scio://contexts/([a-z]+[a-z0-9-]*)$`),
		entityType: model.EntityTypeContext,
	},
}
