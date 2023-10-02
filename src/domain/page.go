package domain

// Page stores information about a page
type Page struct {
	Key          string            `json:"_key,omitempty" example:"123"`
	Information  map[string]string `json:"information,omitempty"`
	Descriptions Descriptions      `json:"descriptions,omitempty"`
	Photos       []Photo           `json:"photos,omitempty"`
	Comments     []Comment         `json:"comments,omitempty"`
}

func (p Page) GetKey() string {
	return p.Key
}

func (p *Page) SetKey(id string) {
	p.Key = id
}
