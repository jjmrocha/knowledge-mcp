package model

import (
	"fmt"
	"time"

	"github.com/jjmrocha/knowledge-mcp/internal/entity"
	"gopkg.in/yaml.v3"
)

type Domain struct {
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

func ParseDomain(content string) (*Domain, error) {
	entityContent, err := entity.ParseContent(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse domain file: %w", err)
	}

	var d Domain
	if err := yaml.Unmarshal([]byte(entityContent.Metadata), &d); err != nil {
		return nil, fmt.Errorf("failed to unmarshal domain metadata: %w", err)
	}

	if d.Entity != EntityTypeDomain {
		return nil, fmt.Errorf("invalid entity type: expected '%s', got '%s'", EntityTypeDomain, d.Entity)
	}

	d.Body = entityContent.Body
	return &d, nil
}

func EncodeDomain(d *Domain) (string, error) {
	entityCopy := *d

	if entityCopy.Tags == nil {
		entityCopy.Tags = []string{}
	}

	if entityCopy.Relations == nil {
		entityCopy.Relations = []RelationRef{}
	}

	metadata, err := yaml.Marshal(&entityCopy)
	if err != nil {
		return "", fmt.Errorf("failed to encode domain: %w", err)
	}

	content := entity.Encode(&entity.EntityContent{
		Metadata: string(metadata),
		Body:     d.Body,
	})

	return content, nil
}
