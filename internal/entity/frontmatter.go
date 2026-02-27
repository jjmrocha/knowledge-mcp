package entity

import (
	"errors"
	"strings"
)

type EntityContent struct {
	Metadata string
	Body     string
}

func ParseContent(content string) (*EntityContent, error) {
	var parsed EntityContent

	if !strings.HasPrefix(content, "---\n") {
		return nil, errors.New("content does not start with YAML frontmatter delimiter")
	}

	rest := content[4:]

	front, tail, found := strings.Cut(rest, "\n---\n")
	if !found {
		return nil, errors.New("missing closing YAML frontmatter delimiter")
	}

	parsed.Metadata = front
	parsed.Body = tail

	return &parsed, nil
}

func Encode(content *EntityContent) string {
	var builder strings.Builder

	builder.WriteString("---\n")
	builder.WriteString(content.Metadata)
	builder.WriteString("\n---\n")
	builder.WriteString(content.Body)

	return builder.String()
}
