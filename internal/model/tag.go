package model

import (
	"fmt"
	"time"

	"github.com/jjmrocha/knowledge-mcp/internal/entity"
	"gopkg.in/yaml.v3"
)

type Tag struct {
	Entity          string    `yaml:"entity"`
	Schema          int       `yaml:"schema"`
	URI             string    `yaml:"uri"`
	Version         int       `yaml:"version"`
	Created         time.Time `yaml:"created"`
	LastUpdate      time.Time `yaml:"last-update"`
	AllowedEntities []string  `yaml:"allowed-entities"`
	Broader         []string  `yaml:"broader"`
	Narrower        []string  `yaml:"narrower"`
	Body            string    `yaml:"-"`
}

func ParseTag(content string) (*Tag, error) {
	entityContent, err := entity.ParseContent(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tag file: %w", err)
	}

	var t Tag
	if err := yaml.Unmarshal([]byte(entityContent.Metadata), &t); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tag metadata: %w", err)
	}

	if t.Entity != EntityTypeTag {
		return nil, fmt.Errorf("invalid entity type: expected '%s', got '%s'", EntityTypeTag, t.Entity)
	}

	t.Body = entityContent.Body
	return &t, nil
}

func EncodeTag(t *Tag) (string, error) {
	entityCopy := *t

	if entityCopy.AllowedEntities == nil {
		entityCopy.AllowedEntities = []string{
			EntityTypeContext,
			EntityTypeDomain,
			EntityTypeConcept,
		}
	}

	if entityCopy.Broader == nil {
		entityCopy.Broader = []string{}
	}

	if entityCopy.Narrower == nil {
		entityCopy.Narrower = []string{}
	}

	metadata, err := yaml.Marshal(&entityCopy)
	if err != nil {
		return "", fmt.Errorf("failed to encode tag: %w", err)
	}

	content := entity.Encode(&entity.EntityContent{
		Metadata: string(metadata),
		Body:     t.Body,
	})

	return content, nil
}
