package storage

import (
	"path/filepath"

	"github.com/jjmrocha/knowledge-mcp/internal/model"
	"github.com/jjmrocha/knowledge-mcp/internal/uri"
)

func FileName(rootDir string, u *uri.URI) string {
	switch {
	case u.Entity == model.EntityTypeTag && u.Context == nil:
		return filepath.Join(rootDir, "tags", u.Slug+".md")

	case u.Entity == model.EntityTypeTag:
		return filepath.Join(rootDir, "contexts", *u.Context, "tags", u.Slug+".md")

	case u.Entity == model.EntityTypeRelation && u.Context == nil:
		return filepath.Join(rootDir, "relations", u.Slug+".md")

	case u.Entity == model.EntityTypeRelation:
		return filepath.Join(rootDir, "contexts", *u.Context, "relations", u.Slug+".md")

	case u.Entity == model.EntityTypeContext:
		return filepath.Join(rootDir, "contexts", u.Slug, "context.md")

	case u.Entity == model.EntityTypeDomain:
		return filepath.Join(rootDir, "contexts", *u.Context, "domains", u.Slug, "domain.md")

	default:
		return filepath.Join(rootDir, "contexts", *u.Context, "domains", *u.Domain, u.Slug+".md")
	}
}

func FileDir(rootDir string, u *uri.URI) string {
	return filepath.Dir(FileName(rootDir, u))
}
