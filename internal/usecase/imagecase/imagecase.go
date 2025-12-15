package imagecase

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type ImageService interface {
	ImageUpload(ctx context.Context, file *multipart.FileHeader) (string, error)
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
	newFileName := fmt.Sprintf("%d%s", time.Now().Unix(), extension)

	dstPath := filepath.Join(s.StoragePath, newFileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	publicPath := filepath.Join(s.PublicPath, newFileName)
	return publicPath, nil

}
