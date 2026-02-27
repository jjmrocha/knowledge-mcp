package v1

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/jjmrocha/knowledge-mcp/internal/model"
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
	entityFile, err := model.ParseEntityFile(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tag file: %w", err)
	}

	var t Tag
	if err := yaml.Unmarshal([]byte(entityFile.Metadata), &t); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tag metadata: %w", err)
	}

	if t.Entity != model.EntityTypeTag {
		return nil, fmt.Errorf("invalid entity type: expected '%s', got '%s'", model.EntityTypeTag, t.Entity)
	}

	t.Body = entityFile.Body
	return &t, nil
}

func EncodeTag(t *Tag) (string, error) {
	copy := *t

	if copy.AllowedEntities == nil {
		copy.AllowedEntities = []string{
			model.EntityTypeContext,
			model.EntityTypeDomain,
			model.EntityTypeConcept,
		}
	}

	if copy.Broader == nil {
		copy.Broader = []string{}
	}

	if copy.Narrower == nil {
		copy.Narrower = []string{}
	}

	metadata, err := yaml.Marshal(&copy)
	if err != nil {
		return "", fmt.Errorf("failed to encode tag: %w", err)
	}

	content := model.EncodeEntityFile(&model.EntityFile{
		Metadata: string(metadata),
		Body:     t.Body,
	})

	return content, nil
}
