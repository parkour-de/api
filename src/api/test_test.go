package api

import (
	"log"
	"os"
	"pkv/api/src/repository/dpv"
	"pkv/api/src/repository/graph"
	"testing"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	Cleanup()
	os.Exit(exitCode)
}

func Cleanup() {
	var err error
	config, err := dpv.NewConfig("../../config.yml")
	if err != nil {
		log.Fatalf("could not initialise config instance: %s", err)
	}
	c, err := graph.Connect(config, true)
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}
	err = graph.DropTestDatabases(c)
}
