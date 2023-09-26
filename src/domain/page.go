package domain

// Page
// @description Stores information about a page
type Page struct {
	ID           string            `json:"_id,omitempty" example:"page/123"`
	Information  map[string]string `json:"information,omitempty"`
	Descriptions Descriptions      `json:"descriptions,omitempty"`
	Photos       []Photo           `json:"photos,omitempty"`
	Comments     []Comment         `json:"comments,omitempty"`
}

func (p Page) GetID() string {
	return p.ID
}

func (p *Page) SetID(id string) {
	p.ID = id
}
