package domain

type TotpConfiguration struct {
	LoginId string `json:"loginId,omitempty" example:"123"`
	Secret  string `json:"secret,omitempty" example:"base32-encoded string"`
	Image   string `json:"image,omitempty" example:"data:image/png;base64,iVBORw0K..."`
}
