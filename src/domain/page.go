package domain

// Page stores information about a page
type Page struct {
	Entity
	Information  map[string]string `json:"information,omitempty"`
	Descriptions Descriptions      `json:"descriptions,omitempty"`
	Photos
	Comments []Comment `json:"comments,omitempty"`
}
