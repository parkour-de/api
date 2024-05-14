package verband

import (
	"reflect"
	"testing"
)

func TestVereineByBundesland(t *testing.T) {
	vereine := []Verein{
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
