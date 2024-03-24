package domain

// User stores information about a user
type User struct {
	Entity
	Name         string            `json:"name,omitempty" example:"John Doe"`
	Type         string            `json:"type,omitempty" example:"person"` // person, team, group, association, company
	Information  map[string]string `json:"information,omitempty"`
	Descriptions Descriptions      `json:"descriptions,omitempty"`
	Photos
	Comments []Comment `json:"comments,omitempty"`
}
