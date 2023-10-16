package user

import (
	"log"
	"pkv/api/src/repository/dpv"
	"testing"
)

func TestSuggest(t *testing.T) {
	config, err := dpv.NewConfig("../../../config.yml")
	if err != nil {
		log.Fatal(err)
	}
	dpv.ConfigInstance = config

	val, err := Suggest()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		val, err = Suggest()
		t.Log("Suggest() = ", val)
	}
}
