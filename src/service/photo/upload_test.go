package photo

import (
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

	s := &Service{}

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
