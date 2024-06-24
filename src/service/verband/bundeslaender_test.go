package verband

import (
	"pkv/api/src/domain/verband"
	"reflect"
	"testing"
)

func TestVereineByBundesland(t *testing.T) {
	vereine := []verband.Verein{
		{Bundesland: "Z", Stadt: "X", Name: "A"},
		{Bundesland: "A", Stadt: "C", Name: "B"},
		{Bundesland: "Z", Stadt: "Y", Name: "A"},
	}

	expected := map[string]int{
		"Z": 2,
		"A": 1,
	}

	actual := aggregateVereineByBundesland(vereine)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("AggregateVereineByBundesland did not give the correct results!\nexpected:\n%#v, \ngot:\n%#v", expected, actual)
	}
}

func TestMitgliederByBundesland(t *testing.T) {
	vereine := []VereinDetail{
		{Bundesland: "Z", Stadt: "X", Name: "A", Mitglieder: 30},
		{Bundesland: "A", Stadt: "C", Name: "B", Mitglieder: 10},
		{Bundesland: "Z", Stadt: "Y", Name: "A", Mitglieder: 20},
	}

	expected := map[string]int{
		"Z": 50,
		"A": 10,
	}

	actual := aggregateMitgliederByBundesland(vereine)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("AggregateMitgliederByBundesland did not give the correct results!\nexpected:\n%#v, \ngot:\n%#v", expected, actual)
	}
}
