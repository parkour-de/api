package domain

type Edge struct {
	ID    string `json:"_id,omitempty"`
	From  string `json:"_from,omitempty"`
	To    string `json:"_to,omitempty"`
	Label string `json:"label,omitempty"`
}
