package domain

// Descriptions
// @description A map of language code to title and text
type Descriptions map[string]Description // per language

type Description struct {
	Title      string `json:"title,omitempty" example:"My Item"`
	Text       string `json:"text,omitempty" example:"Something to describe"`
	Translated bool   `json:"translated,omitempty" example:"false"` // whether this has automatically been translated from another language
}
