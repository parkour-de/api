package verband

import (
	"pkv/api/src/domain/verband"
	"testing"
)

func TestMitmachen(t *testing.T) {
	service := NewService()
	err := service.Mitmachen(verband.MitmachenRequest{
		Name:        "Ricarda Mustermann",
		Email:       "ric@muster.test",
		AG:          "bjoern",
		Kompetenzen: "Sehr gut",
		Fragen:      "Keine",
	})
	if err != nil {
		t.Error(err)
	}
}
