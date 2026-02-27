package v1

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/jjmrocha/knowledge-mcp/internal/model"
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
	entityFile, err := model.ParseEntityFile(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse relation file: %w", err)
	}

	var r RelationType
	if err := yaml.Unmarshal([]byte(entityFile.Metadata), &r); err != nil {
		return nil, fmt.Errorf("failed to unmarshal relation metadata: %w", err)
	}

	if r.Entity != model.EntityTypeRelation {
		return nil, fmt.Errorf("invalid entity type: expected '%s', got '%s'", model.EntityTypeRelation, r.Entity)
	}

	r.Body = entityFile.Body
	return &r, nil
}

func EncodeRelationType(r *RelationType) (string, error) {
	copy := *r

	if copy.AllowedSourceEntities == nil {
		copy.AllowedSourceEntities = []string{
			model.EntityTypeContext,
			model.EntityTypeDomain,
			model.EntityTypeConcept,
		}
	}

	if copy.AllowedTargetEntities == nil {
		copy.AllowedTargetEntities = []string{
			model.EntityTypeContext,
			model.EntityTypeDomain,
			model.EntityTypeConcept,
		}
	}

	metadata, err := yaml.Marshal(&copy)
	if err != nil {
		return "", fmt.Errorf("failed to encode relation type: %w", err)
	}

	content := model.EncodeEntityFile(&model.EntityFile{
		Metadata: string(metadata),
		Body:     r.Body,
	})

	return content, nil
}
