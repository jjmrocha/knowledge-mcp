package v1

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/jjmrocha/knowledge-mcp/internal/model"
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
	entityFile, err := model.ParseEntityFile(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse concept file: %w", err)
	}

	var c Concept
	if err := yaml.Unmarshal([]byte(entityFile.Metadata), &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal concept metadata: %w", err)
	}

	if c.Entity != model.EntityTypeConcept {
		return nil, fmt.Errorf("invalid entity type: expected '%s', got '%s'", model.EntityTypeConcept, c.Entity)
	}

	c.Body = entityFile.Body
	return &c, nil
}

func EncodeConcept(c *Concept) (string, error) {
	copy := *c

	if copy.Tags == nil {
		copy.Tags = []string{}
	}

	if copy.Relations == nil {
		copy.Relations = []RelationRef{}
	}

	if copy.Sources == nil {
		copy.Sources = []Source{}
	}

	metadata, err := yaml.Marshal(&copy)
	if err != nil {
		return "", fmt.Errorf("failed to encode concept: %w", err)
	}

	content := model.EncodeEntityFile(&model.EntityFile{
		Metadata: string(metadata),
		Body:     c.Body,
	})

	return content, nil
}
