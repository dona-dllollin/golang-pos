package imagecase

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/dona-dllollin/belajar-clean-arch/pkgs/logger"
)

type ImageService interface {
	ImageUpload(ctx context.Context, file *multipart.FileHeader) (string, error)
	ImageDelete(ctx context.Context, publicPath string) error
}

type ImageUploadService struct {
	StoragePath string
	PublicPath  string
}

func (s *ImageUploadService) ImageUpload(ctx context.Context, file *multipart.FileHeader) (string, error) {

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	if err := os.MkdirAll(s.StoragePath, 0755); err != nil {
		return "", err
	}

	extension := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)

	storagePath := fmt.Sprintf("%s/%s", s.StoragePath, s.PublicPath)
	dstPath := filepath.Join(storagePath, newFileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	publicPath := path.Join(s.PublicPath, newFileName)
	return publicPath, nil

}

func (s *ImageUploadService) ImageDelete(ctx context.Context, publicPath string) error {
	err := os.Remove(fmt.Sprintf("%s/%s", s.StoragePath, publicPath))
	if err != nil {
		logger.Errorf("Failed to delete image: %v", err)
		return err
	}

	return nil
}
