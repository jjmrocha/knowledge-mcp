package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jjmrocha/knowledge-mcp/internal/model"
	"github.com/jjmrocha/knowledge-mcp/internal/uri"
)

const folderPermissions = 0o755
const filePermissions = 0o644

func ensureParentDirs(rootDir string, u *uri.URI) error {
	return os.MkdirAll(FileDir(rootDir, u), folderPermissions)
}

func SaveFile(rootDir string, u *uri.URI, content []byte) error {
	if err := ensureParentDirs(rootDir, u); err != nil {
		return fmt.Errorf("failed to create parent directories: %w", err)
	}

	fileName := FileName(rootDir, u)
	return os.WriteFile(fileName, content, filePermissions)
}

func ReadFile(rootDir string, u *uri.URI) ([]byte, error) {
	fileName := FileName(rootDir, u)
	return os.ReadFile(fileName) //nolint:gosec
}

func DeleteFile(rootDir string, u *uri.URI) error {
	fileName := FileName(rootDir, u)

	if err := os.Remove(fileName); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	if u.Entity == model.EntityTypeDomain || u.Entity == model.EntityTypeContext {
		parentDir := FileDir(rootDir, u)

		if err := removeDir(parentDir); err != nil {
			return fmt.Errorf("failed to remove parent directory: %w", err)
		}
	}

	return nil
}

func removeDir(dirNme string) error {
	entries, err := os.ReadDir(dirNme)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		entryName := filepath.Join(dirNme, entry.Name())

		if entry.IsDir() {
			if err := removeDir(entryName); err != nil {
				return err
			}

			continue
		}

		if err := os.Remove(entryName); err != nil {
			return fmt.Errorf("failed to remove file: %w", err)
		}
	}

	return os.Remove(dirNme)
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
