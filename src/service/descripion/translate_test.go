package description

import (
	"pkv/api/src/repository/dpv"
	"testing"
)

func TestTranslateDocument(t *testing.T) {
	t.Skip("Translation works, skip test to avoid hitting the DeepL API")
	var err error
	dpv.ConfigInstance, err = dpv.NewConfig("../../../config.yml")
	if err != nil {
		t.Errorf("NewConfig() error = %v", err)
	}
	tests := []struct {
		name     string
		text     string
		srcLang  string
		destLang string
		want     string
	}{
		{
			"EnDe",
			"red",
			"en",
			"de",
			"rot",
		},
		{
			"DeEn",
			"rot",
			"de",
			"en",
			"red",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TranslateDocument(tt.text, tt.srcLang, tt.destLang, nil)
			if err != nil {
				t.Errorf("TranslateDocument(%#v, %#v, %#v, %#v) error = %v", tt.text, tt.srcLang, tt.destLang, nil, err)
			}
			if got != tt.want {
				t.Errorf("TranslateDocument(%#v, %#v, %#v, %#v) got = %#v, want %#v", tt.text, tt.srcLang, tt.destLang, nil, got, tt.want)
			}
		})
	}
}
