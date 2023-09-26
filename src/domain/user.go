package domain

// User
// @description Stores information about a user
type User struct {
	ID           string            `json:"_id,omitempty" example:"user/123"`
	Name         string            `json:"name,omitempty" example:"John Doe"`
	Type         string            `json:"type,omitempty" example:"person"` // person, team, group, association, company
	Information  map[string]string `json:"information,omitempty"`
	Descriptions Descriptions      `json:"descriptions,omitempty"`
	Photos       []Photo           `json:"photos,omitempty"`
	Comments     []Comment         `json:"comments,omitempty"`
}

func (u User) GetID() string {
	return u.ID
}

func (u *User) SetID(id string) {
	u.ID = id
}
