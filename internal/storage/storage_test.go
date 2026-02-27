package storage_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jjmrocha/knowledge-mcp/internal/storage"
	"github.com/jjmrocha/knowledge-mcp/internal/uri"
)

func TestInitRootDirs(t *testing.T) {
	// given
	root := t.TempDir()
	// when
	err := storage.InitRootDirs(root)
	// then
	assert.NoError(t, err)

	for _, dir := range []string{"tags", "relations", "contexts"} {
		info, statErr := os.Stat(filepath.Join(root, dir))
		assert.NoError(t, statErr)
		assert.True(t, info.IsDir())
	}
}

func TestInitRootDirs_Idempotent(t *testing.T) {
	// given
	root := t.TempDir()
	// then
	assert.NoError(t, storage.InitRootDirs(root))
	assert.NoError(t, storage.InitRootDirs(root))
}

func TestEnsureParentDirs(t *testing.T) {
	tests := []struct {
		name        string
		rawURI      string
		expectedDir string
	}{
		{
			name:        "global tag",
			rawURI:      "scio://tags/business-rule",
			expectedDir: "tags",
		},
		{
			name:        "context-scoped tag",
			rawURI:      "scio://contexts/ecommerce/tags/local-tag",
			expectedDir: filepath.Join("contexts", "ecommerce", "tags"),
		},
		{
			name:        "context",
			rawURI:      "scio://contexts/ecommerce",
			expectedDir: filepath.Join("contexts", "ecommerce"),
		},
		{
			name:        "domain",
			rawURI:      "scio://contexts/ecommerce/domains/business-rules",
			expectedDir: filepath.Join("contexts", "ecommerce", "domains", "business-rules"),
		},
		{
			name:        "concept",
			rawURI:      "scio://contexts/ecommerce/domains/business-rules/concepts/discount-calculation",
			expectedDir: filepath.Join("contexts", "ecommerce", "domains", "business-rules"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			u, err := uri.Parse(tc.rawURI)
			assert.NoError(t, err)

			err = storage.EnsureParentDirs(root, u)
			assert.NoError(t, err)

			info, err := os.Stat(filepath.Join(root, tc.expectedDir))
			assert.NoError(t, err)
			assert.True(t, info.IsDir())
		})
	}
}

// buildFindFixture creates:
//
//	root/
//	├── a.txt
//	├── b.md
//	└── sub/
//	    ├── c.txt
//	    └── deep/
//	        └── d.txt
func buildFindFixture(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	writeFile(t, filepath.Join(root, "a.txt"))
	writeFile(t, filepath.Join(root, "b.md"))
	writeFile(t, filepath.Join(root, "sub", "c.txt"))
	writeFile(t, filepath.Join(root, "sub", "deep", "d.txt"))
	return root
}

func writeFile(t *testing.T, path string) {
	t.Helper()
	assert.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	assert.NoError(t, os.WriteFile(path, nil, 0o644))
}

func TestFindFiles_NonRecursive(t *testing.T) {
	// given
	root := buildFindFixture(t)

	var got []string
	// when
	err := storage.FindFiles(root, false, func(filename string) {
		got = append(got, filename)
	})
	// then
	assert.NoError(t, err)
	sort.Strings(got)

	expected := []string{
		filepath.Join(root, "a.txt"),
		filepath.Join(root, "b.md"),
		filepath.Join(root, "sub"),
	}
	sort.Strings(expected)

	assert.Equal(t, expected, got)
}

func TestFindFiles_Recursive(t *testing.T) {
	// given
	root := buildFindFixture(t)

	var got []string
	// when
	err := storage.FindFiles(root, true, func(filename string) {
		got = append(got, filename)
	})
	// then
	assert.NoError(t, err)
	sort.Strings(got)

	expected := []string{
		filepath.Join(root, "a.txt"),
		filepath.Join(root, "b.md"),
		filepath.Join(root, "sub", "c.txt"),
		filepath.Join(root, "sub", "deep", "d.txt"),
	}
	sort.Strings(expected)

	assert.Equal(t, expected, got)
}

func TestFindFiles_Error(t *testing.T) {
	// when
	err := storage.FindFiles("/nonexistent/path", false, func(_ string) {})
	// then
	assert.Error(t, err)
}
