package domain

type Login struct {
	Entity
	Provider string `json:"provider,omitempty" example:"facebook"`
	Subject  string `json:"subject,omitempty" example:"10150000000001234"`
	Enabled  bool   `json:"enabled,omitempty"`
}
