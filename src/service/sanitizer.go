package service

import "github.com/microcosm-cc/bluemonday"

var p *bluemonday.Policy

func SanitizeHTML(input string) string {
	if p == nil {
		p = bluemonday.UGCPolicy().AddTargetBlankToFullyQualifiedLinks(true)
	}
	return p.Sanitize(input)
}
