package domain

type TotpEnableRequest struct {
	LoginId string `json:"loginId,omitempty" example:"123"`
	Code    string `json:"code,omitempty" example:"123456"`
}
