package user

import (
	"context"
	"pkv/api/src/domain"
	"pkv/api/src/repository/t"
	"pkv/api/src/service/description"
	"time"
)

func (s *Service) AddComment(key string, author string, title string, text string, ctx context.Context) error {
	user, err := s.db.Users.Read(key, ctx)
	if err != nil {
		return t.Errorf("read user failed: %w", err)
	}
	text = description.FixTitle(title, text)
	title = description.GetTitle(text)
	render := description.Render([]byte(text))
	if title == "" {
		return t.Errorf("title cannot be empty")
	}
	if len(title) > 100 {
		return t.Errorf("title cannot be longer than 100 characters")
	}
	if text == "" {
		return t.Errorf("text cannot be empty")
	}
	if len(text) > 10000 {
		return t.Errorf("text cannot be longer than 10000 characters")
	}
	comment := domain.Comment{
		Title:   title,
		Text:    text,
		Render:  render,
		Author:  author,
		Created: time.Now(),
	}
	for _, c := range user.Comments {
		if c.Title == title {
			return t.Errorf("comment with same title already exists")
		}
	}
	user.Comments = append(user.Comments, comment)
	if err = s.db.Users.Update(user, ctx); err != nil {
		return t.Errorf("update user failed: %w", err)
	}
	return nil
}

func (s *Service) EditComment(key string, author string, oldTitle string, title string, text string, ctx context.Context) error {
	user, err := s.db.Users.Read(key, ctx)
	if err != nil {
		return t.Errorf("read user failed: %w", err)
	}
	text = description.FixTitle(title, text)
	title = description.GetTitle(text)
	render := description.Render([]byte(text))
	if title == "" {
		return t.Errorf("title cannot be empty")
	}
	if len(title) > 100 {
		return t.Errorf("title cannot be longer than 100 characters")
	}
	if text == "" {
		return t.Errorf("text cannot be empty")
	}
	if len(text) > 10000 {
		return t.Errorf("text cannot be longer than 10000 characters")
	}
	var comment *domain.Comment
	for n, c := range user.Comments {
		if c.Title == oldTitle {
			comment = &user.Comments[n]
		} else if c.Title == title {
			return t.Errorf("comment with same title already exists")
		}
	}
	if comment == nil {
		return t.Errorf("comment not found")
	}
	if comment.Author != author {
		return t.Errorf("not authorized to edit comment")
	}
	comment.Title = title
	comment.Text = text
	comment.Render = render
	if err = s.db.Users.Update(user, ctx); err != nil {
		return t.Errorf("update user failed: %w", err)
	}
	return nil
}

func (s *Service) DeleteComment(key string, author string, title string, ctx context.Context) error {
	user, err := s.db.Users.Read(key, ctx)
	if err != nil {
		return t.Errorf("read user failed: %w", err)
	}
	var newComments []domain.Comment
	for _, c := range user.Comments {
		if c.Title != title {
			newComments = append(newComments, c)
		} else if c.Author != author {
			return t.Errorf("not authorized to delete comment")
		}
	}
	if len(newComments) == len(user.Comments) {
		return t.Errorf("comment not found")
	}
	user.Comments = newComments
	if err = s.db.Users.Update(user, ctx); err != nil {
		return t.Errorf("update user failed: %w", err)
	}
	return nil
}
