package photo

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"pkv/api/src/domain"
	"pkv/api/src/repository/dpv"
	"testing"
)

func setup(t *testing.T, imgDir string, tmpDir string) {
	config, err := dpv.NewConfig("../../../config.yml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	dpv.ConfigInstance = config

	dpv.ConfigInstance.Server.ImgPath = imgDir
	dpv.ConfigInstance.Server.TmpPath = tmpDir
}

func TestService_Update_DontCareIfUnchanged(t *testing.T) {
	imgDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(imgDir)
	tmpDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	setup(t, imgDir, tmpDir)

	s := NewService()
	photos := []domain.Photo{
		{Src: "1234"},
		{Src: "5678"},
	}

	newPhotos, err := s.Update(photos, []string{"1234", "5678"}, context.Background())
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	if len(newPhotos) != 2 {
		t.Errorf("Expected 2 photos, got %v", len(newPhotos))
	}
	if newPhotos[0].Src != "1234" || newPhotos[1].Src != "5678" {
		t.Errorf("Expected photos 1234 and 5678, got %v", newPhotos)
	}
}

func TestService_Update_Reorder(t *testing.T) {
	imgDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(imgDir)
	tmpDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	setup(t, imgDir, tmpDir)

	s := NewService()
	photos := []domain.Photo{
		{Src: "1234"},
		{Src: "5678"},
	}

	newPhotos, err := s.Update(photos, []string{"5678", "1234"}, context.Background())
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	if len(newPhotos) != 2 {
		t.Errorf("Expected 2 photos, got %v", len(newPhotos))
	}
	if newPhotos[0].Src != "5678" || newPhotos[1].Src != "1234" {
		t.Errorf("Expected photos 5678 and 1234, got %v", newPhotos)
	}
}

func TestService_Update_AddPhotos(t *testing.T) {
	imgDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(imgDir)
	tmpDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	setup(t, imgDir, tmpDir)

	s := NewService()
	photos := []domain.Photo{
		{Src: "existing"},
	}

	createJson(t, filepath.Join(dpv.ConfigInstance.Server.TmpPath, "addition.json"), "addition")
	createJson(t, filepath.Join(dpv.ConfigInstance.Server.TmpPath, "ignore.json"), "ignore")

	newPhotos, err := s.Update(photos, []string{"existing", "addition"}, context.Background())
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	if len(newPhotos) != 2 {
		t.Errorf("Expected 2 photos, got %v", len(newPhotos))
	}
	if newPhotos[0].Src != "existing" || newPhotos[1].Src != "addition" {
		t.Errorf("Expected photos existing and addition, got %v", newPhotos)
	}
	if newPhotos[1].W != 42 {
		t.Errorf("Expected photo addition to have width 42, got %v", newPhotos[1].W)
	}

	_, err = s.Update(photos, []string{"existing", "addition"}, context.Background())
	if err == nil {
		t.Errorf("Should have failed as file is now missing, but has not")
	}
}

func TestService_Update_DeletePhotos(t *testing.T) {
	imgDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(imgDir)
	tmpDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	setup(t, imgDir, tmpDir)

	s := NewService()
	photos := []domain.Photo{
		{Src: "existing"},
	}

	createJson(t, filepath.Join(dpv.ConfigInstance.Server.ImgPath, "existing.json"), "existing")

	newPhotos, err := s.Update(photos, []string{}, context.Background())
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	if len(newPhotos) != 0 {
		t.Errorf("Expected 0 photos, got %v", len(newPhotos))
	}

	_, err = s.Update(photos, []string{}, context.Background())
	if err == nil {
		t.Errorf("Should have failed as file is now missing, but has not")
	}
}

func createJson(t *testing.T, path string, src string) {
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create file %s: %v", path, err)
	}
	defer file.Close()

	photo := domain.Photo{Src: src, W: 42}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(photo)

	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
}
