package model

import (
	"fmt"
	"time"

	"github.com/jjmrocha/knowledge-mcp/internal/entity"
	"gopkg.in/yaml.v3"
)

type RelationType struct {
	Entity                string    `yaml:"entity"`
	Schema                int       `yaml:"schema"`
	URI                   string    `yaml:"uri"`
	Version               int       `yaml:"version"`
	Created               time.Time `yaml:"created"`
	LastUpdate            time.Time `yaml:"last-update"`
	InverseOf             string    `yaml:"inverse-of,omitempty"`
	AllowedSourceEntities []string  `yaml:"allowed-source-entities"`
	AllowedTargetEntities []string  `yaml:"allowed-target-entities"`
	Transitive            bool      `yaml:"transitive"`
	Symmetric             bool      `yaml:"symmetric"`
	Body                  string    `yaml:"-"`
}

func ParseRelationType(content string) (*RelationType, error) {
	entityContent, err := entity.ParseContent(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse relation file: %w", err)
	}

	var r RelationType
	if err := yaml.Unmarshal([]byte(entityContent.Metadata), &r); err != nil {
		return nil, fmt.Errorf("failed to unmarshal relation metadata: %w", err)
	}

	if r.Entity != EntityTypeRelation {
		return nil, fmt.Errorf("invalid entity type: expected '%s', got '%s'", EntityTypeRelation, r.Entity)
	}

	r.Body = entityContent.Body
	return &r, nil
}

func EncodeRelationType(r *RelationType) (string, error) {
	entityCopy := *r

	if entityCopy.AllowedSourceEntities == nil {
		entityCopy.AllowedSourceEntities = []string{
			EntityTypeContext,
			EntityTypeDomain,
			EntityTypeConcept,
		}
	}

	if entityCopy.AllowedTargetEntities == nil {
		entityCopy.AllowedTargetEntities = []string{
			EntityTypeContext,
			EntityTypeDomain,
			EntityTypeConcept,
		}
	}

	metadata, err := yaml.Marshal(&entityCopy)
	if err != nil {
		return "", fmt.Errorf("failed to encode relation type: %w", err)
	}

	content := entity.Encode(&entity.EntityContent{
		Metadata: string(metadata),
		Body:     r.Body,
	})

	return content, nil
}
