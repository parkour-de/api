package domain

// Descriptions are maps of language code to title and text
type Descriptions map[string]Description // per language

type Description struct {
	Title      string `json:"title,omitempty" example:"My Item"`
	Text       string `json:"text,omitempty" example:"Something to describe"`
	Render     string `json:"render,omitempty" example:"<p>Something to describe</p>"`
	Translated bool   `json:"translated,omitempty" example:"false"` // whether this has automatically been translated from another language
}
