package domain

type Edge struct {
	Key   string `json:"_key,omitempty"`
	From  string `json:"_from,omitempty"`
	To    string `json:"_to,omitempty"`
	Label string `json:"label,omitempty"`
}
