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

func TestSaveFile(t *testing.T) {
	tests := []struct {
		name         string
		rawURI       string
		expectedFile string
	}{
		{
			name:         "global tag",
			rawURI:       "scio://tags/business-rule",
			expectedFile: filepath.Join("tags", "business-rule.md"),
		},
		{
			name:         "context-scoped tag",
			rawURI:       "scio://contexts/ecommerce/tags/local-tag",
			expectedFile: filepath.Join("contexts", "ecommerce", "tags", "local-tag.md"),
		},
		{
			name:         "context",
			rawURI:       "scio://contexts/ecommerce",
			expectedFile: filepath.Join("contexts", "ecommerce", "context.md"),
		},
		{
			name:         "domain",
			rawURI:       "scio://contexts/ecommerce/domains/business-rules",
			expectedFile: filepath.Join("contexts", "ecommerce", "domains", "business-rules", "domain.md"),
		},
		{
			name:         "concept",
			rawURI:       "scio://contexts/ecommerce/domains/business-rules/concepts/discount-calculation",
			expectedFile: filepath.Join("contexts", "ecommerce", "domains", "business-rules", "discount-calculation.md"),
		},
		{
			name:         "global relation",
			rawURI:       "scio://relations/implements",
			expectedFile: filepath.Join("relations", "implements.md"),
		},
		{
			name:         "context-scoped relation",
			rawURI:       "scio://contexts/ecommerce/relations/owns",
			expectedFile: filepath.Join("contexts", "ecommerce", "relations", "owns.md"),
		},
	}

	content := []byte("hello")

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			u, err := uri.Parse(tc.rawURI)
			assert.NoError(t, err)

			err = storage.SaveFile(root, u, content)
			assert.NoError(t, err)

			got, err := os.ReadFile(filepath.Join(root, tc.expectedFile))
			assert.NoError(t, err)
			assert.Equal(t, content, got)
		})
	}
}

func TestSaveFile_Overwrite(t *testing.T) {
	// given
	root := t.TempDir()
	u := createFile(t, root, "scio://tags/business-rule")
	// when
	err := storage.SaveFile(root, u, []byte("updated"))
	// then
	assert.NoError(t, err)
	got, err := storage.ReadFile(root, u)
	assert.NoError(t, err)
	assert.Equal(t, []byte("updated"), got)
}

func TestReadFile(t *testing.T) {
	// given
	root := t.TempDir()
	u := createFile(t, root, "scio://tags/business-rule")
	// when
	got, err := storage.ReadFile(root, u)
	// then
	assert.NoError(t, err)
	assert.Equal(t, []byte("test file content"), got)
}

func TestReadFile_NotFound(t *testing.T) {
	// given
	root := t.TempDir()
	u, err := uri.Parse("scio://tags/business-rule")
	assert.NoError(t, err)
	// when
	_, err = storage.ReadFile(root, u)
	// then
	assert.Error(t, err)
}

func TestDeleteFile(t *testing.T) {
	// given
	root := t.TempDir()
	u := createFile(t, root, "scio://tags/business-rule")
	// when
	err := storage.DeleteFile(root, u)
	// then
	assert.NoError(t, err)
	_, statErr := os.Stat(filepath.Join(root, "tags", "business-rule.md"))
	assert.True(t, os.IsNotExist(statErr))
}

func TestDeleteFile_NotFound(t *testing.T) {
	// given
	root := t.TempDir()
	u, err := uri.Parse("scio://tags/business-rule")
	assert.NoError(t, err)
	// when
	err = storage.DeleteFile(root, u)
	// then
	assert.Error(t, err)
}

func TestDeleteFile_Context(t *testing.T) {
	// given
	root := t.TempDir()
	u := createFile(t, root, "scio://contexts/ecommerce")
	createFile(t, root, "scio://contexts/ecommerce/domains/business-rules")
	// when
	err := storage.DeleteFile(root, u)
	// then
	assert.NoError(t, err)
	_, statErr := os.Stat(filepath.Join(root, "contexts", "ecommerce"))
	assert.True(t, os.IsNotExist(statErr))
}

func TestDeleteFile_Domain(t *testing.T) {
	// given
	root := t.TempDir()
	u := createFile(t, root, "scio://contexts/ecommerce/domains/business-rules")
	createFile(t, root, "scio://contexts/ecommerce/domains/business-rules/concepts/discount")
	// when
	err := storage.DeleteFile(root, u)
	// then
	assert.NoError(t, err)
	_, statErr := os.Stat(filepath.Join(root, "contexts", "ecommerce", "domains", "business-rules"))
	assert.True(t, os.IsNotExist(statErr))
}

// buildFindFixture creates:
//
//	root/
//	├── contexts/
//	│   └── ecommerce/
//	│       ├── context.md
//	│       └── domains/
//	│           └── rules/
//	│               └── discount.md
//	└── tags/
//	    ├── a.md
//	    └── b.md
func buildFindFixture(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	createFile(t, root, "scio://tags/a")
	createFile(t, root, "scio://tags/b")
	createFile(t, root, "scio://contexts/ecommerce")
	createFile(t, root, "scio://contexts/ecommerce/domains/rules/concepts/discount")
	return root
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
		filepath.Join(root, "contexts"),
		filepath.Join(root, "tags"),
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
		filepath.Join(root, "tags", "a.md"),
		filepath.Join(root, "tags", "b.md"),
		filepath.Join(root, "contexts", "ecommerce", "context.md"),
		filepath.Join(root, "contexts", "ecommerce", "domains", "rules", "discount.md"),
	}
	sort.Strings(expected)

	assert.Equal(t, expected, got)
}

func TestFindFiles_EmptyDir(t *testing.T) {
	// given
	root := t.TempDir()
	var got []string
	// when
	err := storage.FindFiles(root, false, func(filename string) {
		got = append(got, filename)
	})
	// then
	assert.NoError(t, err)
	assert.Empty(t, got)
}

func TestFindFiles_Error(t *testing.T) {
	// when
	err := storage.FindFiles("/nonexistent/path", false, func(_ string) {})
	// then
	assert.Error(t, err)
}

func createFile(t *testing.T, rootDir, fileURI string) *uri.URI {
	t.Helper()
	content := []byte("test file content")
	u, err := uri.Parse(fileURI)
	assert.NoError(t, err)
	err = storage.SaveFile(rootDir, u, content)
	assert.NoError(t, err)
	return u
}
