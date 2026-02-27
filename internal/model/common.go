package model

// RelationRef is a directed edge from a source entity to a target entity
// via a named relation type.
type RelationRef struct {
	Type   string `yaml:"type"`
	Target string `yaml:"target"`
}

// Source is a reference to an external artifact (file, URL, etc.) that
// supports or documents a concept.
type Source struct {
	Type string `yaml:"type"`
	Href string `yaml:"href"`
}
