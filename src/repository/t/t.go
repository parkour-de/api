package t

import (
	"fmt"
	"os"
	"pkv/api/src/repository/dpv"
	"strings"
)

func T(text string) string {
	val, ok := de[text]
	if ok {
		return val
	}
	return text
}

func Errorf(format string, a ...any) error {
	return fmt.Errorf(T(format), a...)
}

var de = map[string]string{}

func LoadDE(config *dpv.Config) error {
	if config == nil {
		return fmt.Errorf("config is not initialized")
	}
	path := config.Path + "strings_de.ini"
	bytes, err := os.ReadFile(path)
	if err != nil {
		wd, _ := os.Getwd()
		return fmt.Errorf("could not load strings_de.ini, looking for %v in %v: %w", path, wd, err)
	}
	for _, s := range strings.Split(string(bytes), "\n") {
		arr := strings.Split(s, "=")
		if len(arr) != 2 {
			return fmt.Errorf("entry contains not exactly one equals sign: %v", s)
		}
		if len(arr[1]) > 0 {
			de[arr[0]] = arr[1]
		}
	}
	return nil
}
