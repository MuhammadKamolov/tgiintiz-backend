package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type FileStorage interface {
	Save(originalName string, src io.Reader) (string, error)
	Delete(filePath string) error
}

type LocalStorage struct {
	uploadDir string
}

func NewLocalStorage(uploadDir string) (*LocalStorage, error) {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("storage: create upload dir: %w", err)
	}
	return &LocalStorage{uploadDir: uploadDir}, nil
}

func (s *LocalStorage) Save(originalName string, src io.Reader) (string, error) {
	ext := filepath.Ext(originalName)
	uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	dstPath := filepath.Join(s.uploadDir, uniqueName)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("storage: create dest file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		_ = os.Remove(dstPath)
		return "", fmt.Errorf("storage: write file: %w", err)
	}

	return "/uploads/" + uniqueName, nil
}

func (s *LocalStorage) Delete(fileURL string) error {
	fileName := filepath.Base(fileURL)
	return os.Remove(filepath.Join(s.uploadDir, fileName))
}
