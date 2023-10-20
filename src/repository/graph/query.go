package graph

import (
	"strings"
)

type QueryBuilder func() (string, map[string]interface{})

func buildUnsetParts(includeSet map[string]struct{}, prefix string) []string {
	var unsetParts []string
	if _, ok := includeSet[prefix+"photos"]; !ok {
		unsetParts = append(unsetParts, `"photos"`)
	}
	if _, ok := includeSet[prefix+"comments"]; !ok {
		unsetParts = append(unsetParts, `"comments"`)
	}
	return unsetParts
}

func appendUnsetPart(list []string, includeSet map[string]struct{}, key string, field string) []string {
	if _, ok := includeSet[key]; !ok {
		list = append(list, `"`+field+`"`)
	}
	return list
}

func buildUnsetString(sectionStr string, unsetParts []string) string {
	if len(unsetParts) > 0 {
		return "UNSET(" + sectionStr + ", " + strings.Join(unsetParts, ", ") + ")"
	}
	return sectionStr
}
