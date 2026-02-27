package model

import (
	"fmt"
	"time"

	"github.com/jjmrocha/knowledge-mcp/internal/entity"
	"gopkg.in/yaml.v3"
)

type Concept struct {
	Entity     string        `yaml:"entity"`
	Schema     int           `yaml:"schema"`
	URI        string        `yaml:"uri"`
	Name       string        `yaml:"name"`
	Version    int           `yaml:"version"`
	Created    time.Time     `yaml:"created"`
	LastUpdate time.Time     `yaml:"last-update"`
	Tags       []string      `yaml:"tags"`
	Relations  []RelationRef `yaml:"relations"`
	Sources    []Source      `yaml:"sources"`
	Body       string        `yaml:"-"`
}

func ParseConcept(content string) (*Concept, error) {
	entityContent, err := entity.ParseContent(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse concept file: %w", err)
	}

	var c Concept
	if err := yaml.Unmarshal([]byte(entityContent.Metadata), &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal concept metadata: %w", err)
	}

	if c.Entity != EntityTypeConcept {
		return nil, fmt.Errorf("invalid entity type: expected '%s', got '%s'", EntityTypeConcept, c.Entity)
	}

	c.Body = entityContent.Body
	return &c, nil
}

func EncodeConcept(c *Concept) (string, error) {
	entityCopy := *c

	if entityCopy.Tags == nil {
		entityCopy.Tags = []string{}
	}

	if entityCopy.Relations == nil {
		entityCopy.Relations = []RelationRef{}
	}

	if entityCopy.Sources == nil {
		entityCopy.Sources = []Source{}
	}

	metadata, err := yaml.Marshal(&entityCopy)
	if err != nil {
		return "", fmt.Errorf("failed to encode concept: %w", err)
	}

	content := entity.Encode(&entity.EntityContent{
		Metadata: string(metadata),
		Body:     c.Body,
	})

	return content, nil
}
