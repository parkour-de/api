package domain

// UserDTO enriches User with some other related nodes such as Page
type UserDTO struct {
	User
	Pages          []Page   `json:"pages,omitempty"`
	PageKeys       []string `json:"pageIds,omitempty"`
	PagePriorities []int    `json:"pagePriorities,omitempty"`
}
