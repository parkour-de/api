package domain

// Photo identifies size and location of an image
type Photo struct {
	Src   string  `json:"src,omitempty" example:"image.jpg"`
	W     int     `json:"w,omitempty" example:"640"`
	H     int     `json:"h,omitempty" example:"480"`
	Lat   float64 `json:"lat,omitempty" example:"54.3243827819444"`
	Lon   float64 `json:"lon,omitempty" example:"10.1457242963889"`
	Color string  `json:"c,omitempty" example:"itV8NsR8q8ibacicQzREZmd3ZUZVmING"`
}
