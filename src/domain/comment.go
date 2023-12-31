package domain

import "time"

// Comment like a guestbook, news feed or blog
type Comment struct {
	Title   string    `json:"title,omitempty" example:"Hey there"`
	Text    string    `json:"text,omitempty" example:"I have something to say here..."`
	Render  string    `json:"render,omitempty" example:"<p>I have something to say here...</p>"`
	Author  string    `json:"author,omitempty" example:"123"`
	Created time.Time `json:"created,omitempty"` // RFC 3339 date
}
