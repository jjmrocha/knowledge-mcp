package storage

import (
	"os"
	"path/filepath"

	"github.com/jjmrocha/knowledge-mcp/internal/uri"
)

const folderPermissions = 0o755

func EnsureParentDirs(rootDir string, u *uri.URI) error {
	return os.MkdirAll(FileDir(rootDir, u), folderPermissions)
}

func InitRootDirs(rootDir string) error {
	for _, dir := range []string{"tags", "relations", "contexts"} {
		if err := os.MkdirAll(filepath.Join(rootDir, dir), folderPermissions); err != nil {
			return err
		}
	}

	return nil
}

func FindFiles(folder string, recursive bool, handler func(filename string)) error {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, e := range entries {
		fullName := filepath.Join(folder, e.Name())

		if recursive && e.IsDir() {
			if err := FindFiles(fullName, recursive, handler); err != nil {
				return err
			}

			continue
		}

		handler(fullName)
	}

	return nil
}
