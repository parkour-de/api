package photo

import (
	"context"
	"fmt"
	"pkv/api/src/domain"
	"pkv/api/src/repository/dpv"
)

// Manage aims to move images between the temporary and the permanent folder
// When a filename is provided, all files matching providedString.[a-z\.]+
// shall be either touched, moved, copied or deleted according to the functions purpose
// example: filename = "1234", the following files would be affected:
// 1234.o.jxl
// 1234.h.jxl
// 1234.json

/*func (s *Service) Add(photos []domain.Photo, additions []string, ctx context.Context) ([]domain.Photo, error) {
	// photos is a slice of domain.Photo{} structs
	// domain.Photo contains Src, which correspond to the string of additions
	// only add photos that weren't in the domain.Photo slice before
	// for each string in the slice call domain.Photo, error = s.ReadPhoto(string, dpv.ConfigInstance.Server.TmpPath, ctx)
	// then error = s.MakePermanent(string, ctx)
	// finally return the new slice with the added photos

	// Here's the challenging thing: We have to undo the changes by calling MakeTemporary if an error happens after calling any MakePermanent...
	var newPhotos []domain.Photo
	for _, addition := range additions {
		photo, err := s.ReadPhoto(addition, dpv.ConfigInstance.Server.TmpPath, ctx)
		if err != nil {
			return nil, fmt.Errorf("could not read photo information: %w", err)
		}
		newPhotos = append(newPhotos, photo)
	}
	var permanentPhotos []string
	for _, addition := range newPhotos {
		if err := s.MakePermanent(addition.Src, ctx); err != nil {
			err := fmt.Errorf("could not make photo %v permanent: %w", addition.Src, err)
			for _, p := range permanentPhotos {
				if err2 := s.MakeTemporary(p, ctx); err != nil {
					err = fmt.Errorf("%w; reverting %v failed: %v", err, p, err2.Error())
				}
			}
			return nil, err
		}
		permanentPhotos = append(permanentPhotos, addition.Src)
	}
	return append(photos, newPhotos...), nil
}

func (s *Service) Remove(photos []domain.Photo, removals []string, ctx context.Context) ([]domain.Photo, error) {
	// similar to above,
	// but as we only want to remove items with matching Src from the slice, we can skip the ReadPhoto part
	// just don't forget to call error = s.MakeTemporary(string, ctx)

	// Here's the challenging thing: We have to undo the changes by calling MakePermanent if an error happens after calling any MakeTemporary...
	var newPhotos []domain.Photo
	var removedPhotos []string
	for _, photo := range photos {
		shouldRemove := false
		for _, removal := range removals {
			if photo.Src == removal {
				shouldRemove = true
				if err := s.MakeTemporary(removal, ctx); err != nil {
					err := fmt.Errorf("could not make photo %v temporary: %w", removal, err)
					for _, p := range removedPhotos {
						if err2 := s.MakePermanent(p, ctx); err != nil {
							err = fmt.Errorf("%w; reverting %v failed: %v", err, p, err2.Error())
						}
					}
					return nil, err
				}
				removedPhotos = append(removedPhotos, removal)
				break
			}
		}

		if !shouldRemove {
			newPhotos = append(newPhotos, photo)
		}
	}

	return newPhotos, nil
}

func (s *Service) Reorder(photos []domain.Photo, order []string, ctx context.Context) ([]domain.Photo, error) {
	// return the photos slice in the same order as the order slice, by comparing the domain.Photo.Src values with the order slice
	var reorderedPhotos []domain.Photo
	for _, src := range order {
		for _, photo := range photos {
			if photo.Src == src {
				reorderedPhotos = append(reorderedPhotos, photo)
				break
			}
		}
	}

	return reorderedPhotos, nil
}*/

func (s *Service) Update(photos []domain.Photo, files []string, ctx context.Context) ([]domain.Photo, error) {
	// if it is possible to do all of the above in a single function, give it a try!
	if hasDuplicates(files) {
		return photos, fmt.Errorf("input slice contains duplicates")
	}

	var updatedPhotos []domain.Photo
	var addedPhotos []domain.Photo
	var removedPhotos []string

	for _, file := range files {
		found := false
		for _, photo := range photos {
			if photo.Src == file {
				found = true
				updatedPhotos = append(updatedPhotos, photo)
				break
			}
		}

		if !found {
			photo, err := s.ReadPhoto(file, dpv.ConfigInstance.Server.TmpPath, ctx)
			if err != nil {
				return photos, s.Undo(addedPhotos, removedPhotos, ctx, fmt.Errorf("could not read photo information for %v: %w", file, err))
			}
			if err := s.MakePermanent(file, ctx); err != nil {
				return photos, s.Undo(addedPhotos, removedPhotos, ctx, fmt.Errorf("could not make photo %v permanent: %w", photo.Src, err))
			}
			addedPhotos = append(addedPhotos, photo)
		}
	}

	for _, photo := range photos {
		found := false
		for _, file := range files {
			if photo.Src == file {
				found = true
				break
			}
		}

		if !found {
			if err := s.MakeTemporary(photo.Src, ctx); err != nil {
				return photos, s.Undo(addedPhotos, removedPhotos, ctx, fmt.Errorf("could not make photo %v temporary: %w", photo.Src, err))
			}
			removedPhotos = append(removedPhotos, photo.Src)
		}
	}

	updatedPhotos = append(updatedPhotos, addedPhotos...)
	return updatedPhotos, nil
}

func (s *Service) Undo(addedPhotos []domain.Photo, removedPhotos []string, ctx context.Context, err error) error {
	for _, p := range addedPhotos {
		if err2 := s.MakeTemporary(p.Src, ctx); err != nil {
			err = fmt.Errorf("%w; reverting %v to Temporary failed: %v", err, p, err2.Error())
		}
	}
	for _, p := range removedPhotos {
		if err2 := s.MakePermanent(p, ctx); err != nil {
			err = fmt.Errorf("%w; reverting %v to Permanent failed: %v", err, p, err2.Error())
		}
	}
	return err
}

func hasDuplicates(slice []string) bool {
	seen := make(map[string]bool)
	for _, s := range slice {
		if seen[s] {
			return true
		}
		seen[s] = true
	}
	return false
}
