package imagecase

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImageUploadService_Upload(t *testing.T) {
	tmpDir := t.TempDir()

	service := &ImageUploadService{
		StoragePath: tmpDir,
	}

	// fake multipart file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("image", "test.jpg")
	part.Write([]byte("fake image content"))
	writer.Close()

	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ParseMultipartForm(10 << 20)

	file, header, err := req.FormFile("image")
	require.NoError(t, err)
	defer file.Close()

	path, err := service.ImageUpload(context.Background(), header)
	require.NoError(t, err)

	// assertions
	assert.FileExists(t, path)
}
