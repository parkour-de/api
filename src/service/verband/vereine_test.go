package verband

import (
	"reflect"
	"testing"
)

func TestFindByQuestionId(t *testing.T) {
	answers := answerList{
		{QuestionId: 16, Text: "Answer 16"},
		{QuestionId: 17, Text: "Answer 17"},
	}

	foundAnswer := answers.findByQuestionId(16)
	if foundAnswer.Text != "Answer 16" {
		t.Errorf("FindByQuestionId failed, expected 'Answer 16' but got %s", foundAnswer.Text)
	}

	foundAnswer = answers.findByQuestionId(18)
	if foundAnswer.Text != "" {
		t.Errorf("FindByQuestionId failed, expected empty string but got %s", foundAnswer.Text)
	}
}

func TestSortingVereine(t *testing.T) {
	vereine := []Verein{
		{Bundesland: "Z", Stadt: "X", Name: "A"},
		{Bundesland: "A", Stadt: "C", Name: "B"},
		{Bundesland: "Z", Stadt: "Y", Name: "A"},
	}

	NewService().sortVereine(vereine)

	// Add assertions to check if the slice is sorted correctly
	if vereine[0].Bundesland != "A" {
		t.Errorf("Sorting failed, expected Bundesland 'A' but got %s", vereine[0].Bundesland)
	}
	if vereine[1].Bundesland != "Z" {
		t.Errorf("Sorting failed, expected Bundesland 'Z' but got %s", vereine[1].Bundesland)
	}
	if vereine[2].Bundesland != "Z" {
		t.Errorf("Sorting failed, expected Bundesland 'Z' but got %s", vereine[2].Bundesland)
	}
	if vereine[0].Stadt != "C" {
		t.Errorf("Sorting failed, expected Stadt 'C' but got %s", vereine[0].Stadt)
	}
	if vereine[1].Stadt != "X" {
		t.Errorf("Sorting failed, expected Stadt 'X' but got %s", vereine[1].Stadt)
	}
	if vereine[2].Stadt != "Y" {
		t.Errorf("Sorting failed, expected Stadt 'Y' but got %s", vereine[2].Stadt)
	}
	if vereine[0].Name != "B" {
		t.Errorf("Sorting failed, expected Name 'B' but got %s", vereine[0].Name)
	}
	if vereine[1].Name != "A" {
		t.Errorf("Sorting failed, expected Name 'A' but got %s", vereine[1].Name)
	}
	if vereine[2].Name != "A" {
		t.Errorf("Sorting failed, expected Name 'A' but got %s", vereine[2].Name)
	}
}

func TestExtractVereineList(t *testing.T) {
	sampleResponse := nextcloudResponse{
		OCS: ocs{
			Data: data{
				Submissions: []submission{
					{
						Answers: []answer{
							{QuestionId: 16, Text: "Ja"},
							{QuestionId: 17, Text: "  Bundesland1  "},
							{QuestionId: 12, Text: "  Stadt1  "},
							{QuestionId: 13, Text: "  Name1  "},
							{QuestionId: 6, Text: "  Webseite1.de  "},
						},
					},
					{
						Answers: []answer{
							{QuestionId: 16, Text: "Nein"},
							{QuestionId: 17, Text: "Bundesland2"},
							{QuestionId: 12, Text: "Stadt2"},
							{QuestionId: 13, Text: "Name2"},
							{QuestionId: 6, Text: "Webseite2"},
						},
					},
					{
						Answers: []answer{
							{QuestionId: 16, Text: "Ja"},
							{QuestionId: 13, Text: "Name3"},
							{QuestionId: 6, Text: "ftp://ftp.example.com"},
						},
					},
				},
			},
		},
	}

	extractedVereine := NewService().ExtractVereineList(sampleResponse)

	expectedVereine := []Verein{
		{
			Bundesland: "Bundesland1",
			Stadt:      "Stadt1",
			Name:       "Name1",
			Webseite:   "https://Webseite1.de/",
		},
		{
			Bundesland: "",
			Stadt:      "",
			Name:       "Name3",
			Webseite:   "ftp://ftp.example.com/",
		},
	}

	if !reflect.DeepEqual(extractedVereine, expectedVereine) {
		t.Errorf("ExtractVereineList did not return the expected result:\nexpected:\n%#v, \ngot:\n%#v", expectedVereine, extractedVereine)
	}
}
