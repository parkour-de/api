package verband

type MitmachenRequest struct {
	Name        string `json:"name" example:"Erika Mustermann"`
	Email       string `json:"email" example:"erika.mustermann@mail.de"`
	AG          string `json:"ag" example:"oeffentlichkeit"`
	Kompetenzen string `json:"kompetenzen" example:"Ich habe 10 Jahre lang als Nachrichtensprecher gearbeitet."`
	Fragen      string `json:"fragen" example:"Gibt es auch Live-Meetings?"`
	Altcha      string `json:"altcha" example:"BASE64EncodedStringWithASolvedCaptcha"`
}
