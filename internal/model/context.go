package model

import (
	"fmt"
	"time"

	"github.com/jjmrocha/knowledge-mcp/internal/entity"
	"gopkg.in/yaml.v3"
)

type Context struct {
	Entity     string        `yaml:"entity"`
	Schema     int           `yaml:"schema"`
	URI        string        `yaml:"uri"`
	Name       string        `yaml:"name"`
	Version    int           `yaml:"version"`
	Created    time.Time     `yaml:"created"`
	LastUpdate time.Time     `yaml:"last-update"`
	Tags       []string      `yaml:"tags"`
	Relations  []RelationRef `yaml:"relations"`
	Body       string        `yaml:"-"`
}

func ParseContext(content string) (*Context, error) {
	entityContent, err := entity.ParseContent(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse context file: %w", err)
	}

	var c Context
	if err := yaml.Unmarshal([]byte(entityContent.Metadata), &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal context metadata: %w", err)
	}

	if c.Entity != EntityTypeContext {
		return nil, fmt.Errorf("invalid entity type: expected '%s', got '%s'", EntityTypeContext, c.Entity)
	}

	c.Body = entityContent.Body
	return &c, nil
}

func EncodeContext(c *Context) (string, error) {
	entityCopy := *c

	if entityCopy.Tags == nil {
		entityCopy.Tags = []string{}
	}

	if entityCopy.Relations == nil {
		entityCopy.Relations = []RelationRef{}
	}

	metadata, err := yaml.Marshal(&entityCopy)
	if err != nil {
		return "", fmt.Errorf("failed to encode context: %w", err)
	}

	content := entity.Encode(&entity.EntityContent{
		Metadata: string(metadata),
		Body:     c.Body,
	})

	return content, nil
}
