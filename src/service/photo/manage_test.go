package photo

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestService_Clean(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(tempDir, "file2.txt")
	file3 := filepath.Join(tempDir, "file3.txt")

	now := time.Now()
	oldTime := now.Add(-4 * 24 * time.Hour) // 4 days ago
	recentTime := now.Add(-1 * time.Hour)   // 1 hour ago

	createFile(t, file1, oldTime)
	createFile(t, file2, recentTime)
	createFile(t, file3, oldTime)

	s := NewService()

	err = s.Clean(tempDir, 72*time.Hour)
	if err != nil {
		t.Errorf("Clean failed: %v", err)
	}

	_, err = os.Stat(file1)
	if !os.IsNotExist(err) {
		t.Errorf("file1.txt should have been removed")
	}

	_, err = os.Stat(file2)
	if err != nil {
		t.Errorf("file2.txt should still exist: %v", err)
	}

	_, err = os.Stat(file3)
	if !os.IsNotExist(err) {
		t.Errorf("file3.txt should have been removed")
	}
}

func TestService_Touch(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	file1 := filepath.Join(tempDir, "sampleFile.o.jxl")
	file2 := filepath.Join(tempDir, "sampleFile.json")
	file3 := filepath.Join(tempDir, "different.json")

	now := time.Now()
	oldTime := now.Add(-4 * 24 * time.Hour) // 4 days ago
	recentTime := now.Add(-1 * time.Hour)   // 1 hour ago
	almostNow := now.Add(-5 * time.Minute)  // 5 minutes ago

	createFile(t, file1, oldTime)
	createFile(t, file2, recentTime)
	createFile(t, file3, oldTime)

	s := NewService()

	err = s.Touch("sampleFile", tempDir, context.Background())
	if err != nil {
		t.Errorf("Touch failed: %v", err)
	}

	stat, err := os.Stat(file1)
	if err != nil || stat.ModTime().Before(almostNow) {
		t.Errorf("file1.txt should have been touched")
	}

	stat, err = os.Stat(file2)
	if err != nil || stat.ModTime().Before(almostNow) {
		t.Errorf("file2.txt should have been touched")
	}

	stat, err = os.Stat(file3)
	if err != nil || stat.ModTime().After(almostNow) {
		t.Errorf("file3.txt should not have been touched")
	}
}

func TestService_Touch_Error(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	s := NewService()

	err = s.Touch("sampleFile", tempDir, context.Background())
	if err == nil {
		t.Errorf("Touch should have failed")
	}
}

func TestService_Move(t *testing.T) {
	srcDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(srcDir)

	tarDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tarDir)

	file1 := filepath.Join(srcDir, "sampleFile.o.jxl")
	file2 := filepath.Join(srcDir, "sampleFile.json")
	file3 := filepath.Join(srcDir, "different.json")

	createFile(t, file1, time.Now())
	createFile(t, file2, time.Now())
	createFile(t, file3, time.Now())

	s := NewService()

	err = s.Move("sampleFile", srcDir, tarDir, context.Background())
	if err != nil {
		t.Errorf("Move failed: %v", err)
	}

	_, err = os.Stat(filepath.Join(tarDir, "sampleFile.o.jxl"))
	if err != nil {
		t.Errorf("sampleFile.o.jxl should have been moved")
	}

	_, err = os.Stat(filepath.Join(tarDir, "sampleFile.json"))
	if err != nil {
		t.Errorf("sampleFile.json should have been moved")
	}

	_, err = os.Stat(filepath.Join(tarDir, "different.json"))
	if !os.IsNotExist(err) {
		t.Errorf("different.json should not have been moved")
	}

	_, err = os.Stat(filepath.Join(srcDir, "sampleFile.o.jxl"))
	if !os.IsNotExist(err) {
		t.Errorf("sampleFile.o.jxl should be gone")
	}

	_, err = os.Stat(filepath.Join(srcDir, "sampleFile.json"))
	if !os.IsNotExist(err) {
		t.Errorf("sampleFile.json should be gone")
	}

	_, err = os.Stat(filepath.Join(srcDir, "different.json"))
	if err != nil {
		t.Errorf("different.json should not be gone")
	}
}

func TestService_Move_Error(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	s := NewService()

	err = s.Move("sampleFile", tempDir, tempDir, context.Background())
	if err == nil {
		t.Errorf("Move should have failed")
	}
}

func TestService_Copy(t *testing.T) {
	srcDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(srcDir)

	tarDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tarDir)

	file1 := filepath.Join(srcDir, "sampleFile.o.jxl")
	file2 := filepath.Join(srcDir, "sampleFile.json")
	file3 := filepath.Join(srcDir, "different.json")

	createFile(t, file1, time.Now())
	createFile(t, file2, time.Now())
	createFile(t, file3, time.Now())

	// put simple string into file2:

	err = os.WriteFile(file2, []byte(`{"src": "sampleFile", "w": 42}`), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	s := NewService()

	photo, err := s.Copy("sampleFile", srcDir, tarDir, context.Background())
	if err != nil {
		t.Errorf("Copy failed: %v", err)
	}
	if photo.Src == "sampleFile" {
		t.Errorf("photo.Src should have changed")
	}
	if photo.W != 42 {
		t.Errorf("photo.W should be 42")
	}

	_, err = os.Stat(filepath.Join(tarDir, photo.Src+".o.jxl"))
	if err != nil {
		t.Errorf("sampleFile.o.jxl should have been copied")
	}

	_, err = os.Stat(filepath.Join(tarDir, photo.Src+".json"))
	if err != nil {
		t.Errorf("sampleFile.json should have been copied")
	}

	_, err = os.Stat(filepath.Join(tarDir, "different.json"))
	if !os.IsNotExist(err) {
		t.Errorf("different.json should not have been copied")
	}

	_, err = os.Stat(filepath.Join(srcDir, "sampleFile.o.jxl"))
	if err != nil {
		t.Errorf("sampleFile.o.jxl should stay")
	}

	_, err = os.Stat(filepath.Join(srcDir, "sampleFile.json"))
	if err != nil {
		t.Errorf("sampleFile.json should stay")
	}
}

func TestService_Copy_Error(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dpv-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	s := NewService()

	_, err = s.Copy("sampleFile", tempDir, tempDir, context.Background())
	if err == nil {
		t.Errorf("Copy should have failed")
	}
}

func createFile(t *testing.T, path string, modTime time.Time) {
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create file %s: %v", path, err)
	}
	file.Close()

	err = os.Chtimes(path, modTime, modTime)
	if err != nil {
		t.Fatalf("Failed to set ModTime for file %s: %v", path, err)
	}
}
