package domain

// Photo
// @description Identify size and location of an image
type Photo struct {
	Src string `json:"src,omitempty" example:"image.jpg"`
	W   int    `json:"w,omitempty" example:"640"`
	H   int    `json:"h,omitempty" example:"480"`
}
