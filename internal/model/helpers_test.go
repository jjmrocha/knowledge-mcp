package model_test

// entityContent wraps meta and body in YAML frontmatter delimiters,
// producing a valid entity file string for use in tests.
func entityContent(meta, body string) string {
	return "---\n" + meta + "\n---\n" + body
}
