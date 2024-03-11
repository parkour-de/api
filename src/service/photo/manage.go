package photo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"pkv/api/src/domain"
	"pkv/api/src/repository/dpv"
	"regexp"
	"strings"
	"time"
)

// Manage aims to move images between the temporary and the permanent folder
// When a filename is provided, all files matching providedString.[a-z\.]+
// shall be either touched, moved, copied or deleted according to the functions purpose
// example: filename = "1234", the following files would be affected:
// 1234.o.jxl
// 1234.h.jxl
// 1234.json

func (s *Service) ReadPhoto(filename string, path string, ctx context.Context) (domain.Photo, error) {
	if matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]{8,}$", filename); !matched {
		return domain.Photo{}, errors.New("readPhoto: valid filenames can only contain the characters a-z, A-Z, 0-9, _, and -")
	}

	jsonFile, err := os.ReadFile(filepath.Join(path, filename+".json"))
	if err != nil {
		return domain.Photo{}, fmt.Errorf("readPhoto: could not read json file: %w", err)
	}

	var photo domain.Photo
	err = json.Unmarshal(jsonFile, &photo)
	if err != nil {
		return domain.Photo{}, fmt.Errorf("readPhoto: could not decode json file: %w", err)
	}

	return photo, nil
}

// Touch is used to change the filemtime of these files to the current time
func (s *Service) Touch(filename string, path string, ctx context.Context) error {
	if matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]{8,}$", filename); !matched {
		return errors.New("touch: valid filenames can only contain the characters a-z, A-Z, 0-9, _, and -")
	}
	pattern := filepath.Join(path, filename+".*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return errors.New("touch: no matching files found")
	}
	for _, match := range matches {
		err := os.Chtimes(match, time.Now(), time.Now())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) MakePermanent(filename string, ctx context.Context) error {
	if err := s.Move(filename, dpv.ConfigInstance.Server.TmpPath, dpv.ConfigInstance.Server.ImgPath, ctx); err != nil {
		return fmt.Errorf("could not move file: %w", err)
	}
	return nil
}

func (s *Service) MakeTemporary(filename string, ctx context.Context) error {
	// touch first, to avoid immediate deletion of stale files
	if err := s.Touch(filename, dpv.ConfigInstance.Server.ImgPath, ctx); err != nil {
		return fmt.Errorf("could not touch file: %w", err)
	}
	if err := s.Move(filename, dpv.ConfigInstance.Server.ImgPath, dpv.ConfigInstance.Server.TmpPath, ctx); err != nil {
		return fmt.Errorf("could not move file: %w", err)
	}
	return nil
}

func (s *Service) MakeClone(filename string, ctx context.Context) (domain.Photo, error) {
	photo, err := s.Copy(filename, dpv.ConfigInstance.Server.ImgPath, dpv.ConfigInstance.Server.TmpPath, ctx)
	if err != nil {
		return domain.Photo{}, fmt.Errorf("could not clone file: %w", err)
	}
	return photo, nil
}

func (s *Service) Move(filename string, fromPath string, toPath string, ctx context.Context) error {
	if matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]{8,}$", filename); !matched {
		return errors.New("move: valid filenames can only contain the characters a-z, A-Z, 0-9, _, and -")
	}
	fromPattern := filepath.Join(fromPath, filename+".*")

	matches, err := filepath.Glob(fromPattern)
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return errors.New("move: no matching files found")
	}
	for _, fromMatch := range matches {
		toMatch := filepath.Join(toPath, filepath.Base(fromMatch))
		err := os.Rename(fromMatch, toMatch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Copy(filename string, fromPath string, toPath string, ctx context.Context) (domain.Photo, error) {
	if matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]{8,}$", filename); !matched {
		return domain.Photo{}, errors.New("copy: valid filenames can only contain the characters a-z, A-Z, 0-9, _, and -")
	}
	fromPattern := filepath.Join(fromPath, filename+".*")

	matches, err := filepath.Glob(fromPattern)
	if err != nil {
		return domain.Photo{}, err
	}
	if len(matches) == 0 {
		return domain.Photo{}, errors.New("copy: no matching files found")
	}

	newFilename := RandomString()
	var photo domain.Photo

	for _, fromMatch := range matches {
		ext := strings.TrimPrefix(filepath.Base(fromMatch), filename)
		if ext == ".json" {
			photoBytes, err := os.ReadFile(fromMatch)
			if err != nil {
				return domain.Photo{}, fmt.Errorf("copy: could not open json file %s%s: %w", filename, ext, err)
			}
			err = json.Unmarshal(photoBytes, &photo)
			if err != nil {
				return domain.Photo{}, fmt.Errorf("copy: could not decode json file %s%s: %w", filename, ext, err)
			}
			photo.Src = newFilename
		} else {
			fromFile, err := os.Open(fromMatch)
			if err != nil {
				return domain.Photo{}, fmt.Errorf("copy: failed reading file %s%s: %w", filename, ext, err)
			}
			defer fromFile.Close()

			toMatch := filepath.Join(toPath, newFilename+ext)
			toFile, err := os.Create(toMatch)
			if err != nil {
				return domain.Photo{}, fmt.Errorf("copy: failed creating file %s%s: %w", filename, ext, err)
			}
			defer toFile.Close()

			_, err = io.Copy(toFile, fromFile)
			if err != nil {
				return domain.Photo{}, fmt.Errorf("copy: failed copying contents of file %s%s: %w", filename, ext, err)
			}
		}
	}

	if photo.Src == "" {
		return domain.Photo{}, fmt.Errorf("copy: no json file found")
	}

	photoBytes, err := json.Marshal(photo)
	if err != nil {
		return domain.Photo{}, fmt.Errorf("copy: could not encode json file %s: %w", filename, err)
	}

	jsonFilePath := filepath.Join(toPath, newFilename+".json")
	err = os.WriteFile(jsonFilePath, photoBytes, 0644)
	if err != nil {
		return domain.Photo{}, fmt.Errorf("copy: could not save json file %s: %w", filename, err)
	}

	return photo, nil
}

func (s *Service) Clean(folderPath string, maxAge time.Duration) error {
	now := time.Now()
	cutoff := now.Add(-maxAge)

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.ModTime().Before(cutoff) {
			err := os.Remove(path)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
