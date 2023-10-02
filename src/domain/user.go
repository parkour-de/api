package domain

// User stores information about a user
type User struct {
	Key          string            `json:"_key,omitempty" example:"123"`
	Name         string            `json:"name,omitempty" example:"John Doe"`
	Type         string            `json:"type,omitempty" example:"person"` // person, team, group, association, company
	Information  map[string]string `json:"information,omitempty"`
	Descriptions Descriptions      `json:"descriptions,omitempty"`
	Photos       []Photo           `json:"photos,omitempty"`
	Comments     []Comment         `json:"comments,omitempty"`
}

func (u User) GetKey() string {
	return u.Key
}

func (u *User) SetKey(id string) {
	u.Key = id
}
