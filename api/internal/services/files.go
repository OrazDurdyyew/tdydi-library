package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
)

type FileStorage struct {
	root string
}

func NewFileStorage(root string) *FileStorage {
	return &FileStorage{root: root}
}

func (fs *FileStorage) AbsPath(relativePath string) string {
	return filepath.Join(fs.root, relativePath)
}

func (fs *FileStorage) BookDir(bookID int64) string {
	return filepath.Join("books", fmt.Sprint(bookID))
}

func (fs *FileStorage) EnsureDir(relativePath string) error {
	return os.MkdirAll(filepath.Join(fs.root, relativePath), 0755)
}

func (fs *FileStorage) Remove(relativePath string) error {
	if relativePath == "" {
		return nil
	}
	return os.Remove(filepath.Join(fs.root, relativePath))
}

// SaveUpload сохраняет multipart.FileHeader по относительному пути внутри storage root.
func (fs *FileStorage) SaveUpload(fh *multipart.FileHeader, relativePath string) error {
	absPath := filepath.Join(fs.root, relativePath)
	if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
		return err
	}
	src, err := fh.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(absPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// GenerateCoverFromPDF извлекает первую страницу PDF в JPEG через pdftoppm.
// Возвращает относительный путь до созданной обложки (от storage root).
// pdftoppm должен быть установлен в системе (пакет poppler-utils).
func (fs *FileStorage) GenerateCoverFromPDF(pdfRelativePath, outputDir string) (string, error) {
	pdfPath := filepath.Join(fs.root, pdfRelativePath)
	outputAbs := filepath.Join(fs.root, outputDir, "cover")

	cmd := exec.Command("pdftoppm",
		"-jpeg", "-f", "1", "-l", "1", "-r", "150", "-singlefile",
		pdfPath, outputAbs)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pdftoppm: %w", err)
	}

	return filepath.Join(outputDir, "cover.jpg"), nil
}
