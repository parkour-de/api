package user

import (
	"context"
	"fmt"
	"pkv/api/src/domain"
	"time"
)

func (s *Service) AddComment(key string, author string, title string, text string, ctx context.Context) error {
	user, err := s.db.Users.Read(key, ctx)
	if err != nil {
		return fmt.Errorf("read user failed: %w", err)
	}
	comment := domain.Comment{
		Title:   title,
		Text:    text,
		Author:  author,
		Created: time.Now(),
	}
	if title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if len(title) > 100 {
		return fmt.Errorf("title cannot be longer than 100 characters")
	}
	if text == "" {
		return fmt.Errorf("text cannot be empty")
	}
	if len(text) > 10000 {
		return fmt.Errorf("text cannot be longer than 10000 characters")
	}
	for _, c := range user.Comments {
		if c.Title == title {
			return fmt.Errorf("comment with same title already exists")
		}
	}
	user.Comments = append(user.Comments, comment)
	if err = s.db.Users.Update(user, ctx); err != nil {
		return fmt.Errorf("update user failed: %w", err)
	}
	return nil
}

func (s *Service) EditComment(key string, author string, oldTitle string, title string, text string, ctx context.Context) error {
	user, err := s.db.Users.Read(key, ctx)
	if err != nil {
		return fmt.Errorf("read user failed: %w", err)
	}
	if title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if len(title) > 100 {
		return fmt.Errorf("title cannot be longer than 100 characters")
	}
	if text == "" {
		return fmt.Errorf("text cannot be empty")
	}
	if len(text) > 10000 {
		return fmt.Errorf("text cannot be longer than 10000 characters")
	}
	var comment *domain.Comment
	for n, c := range user.Comments {
		if c.Title == oldTitle {
			comment = &user.Comments[n]
		} else if c.Title == title {
			return fmt.Errorf("comment with same title already exists")
		}
	}
	if comment == nil {
		return fmt.Errorf("comment not found")
	}
	if comment.Author != author {
		return fmt.Errorf("not authorized to edit comment")
	}
	comment.Title = title
	comment.Text = text
	if err = s.db.Users.Update(user, ctx); err != nil {
		return fmt.Errorf("update user failed: %w", err)
	}
	return nil
}

func (s *Service) DeleteComment(key string, author string, title string, ctx context.Context) error {
	user, err := s.db.Users.Read(key, ctx)
	if err != nil {
		return fmt.Errorf("read user failed: %w", err)
	}
	var newComments []domain.Comment
	for _, c := range user.Comments {
		if c.Title != title {
			newComments = append(newComments, c)
		} else if c.Author != author {
			return fmt.Errorf("not authorized to delete comment")
		}
	}
	if len(newComments) == len(user.Comments) {
		return fmt.Errorf("comment not found")
	}
	user.Comments = newComments
	if err = s.db.Users.Update(user, ctx); err != nil {
		return fmt.Errorf("update user failed: %w", err)
	}
	return nil
}
