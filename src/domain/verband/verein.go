package verband

import "sort"

// Verein contains information about a verein that wants to be published in a public list
type Verein struct {
	Bundesland string `json:"bundesland" example:"Baden-Württemberg"`
	Stadt      string `json:"stadt" example:"Karlsruhe"`
	Name       string `json:"name" example:"Vereinsname"`
	Webseite   string `json:"webseite" example:"https://verein.karlsruhe.de/"`
}

func (v Verein) SortKey() string {
	return v.Bundesland + "|" + v.Stadt + "|" + v.Name
}

func SortVereine(vereine []Verein) {
	sort.Slice(vereine, func(i, j int) bool {
		return vereine[i].SortKey() < vereine[j].SortKey()
	})
}
