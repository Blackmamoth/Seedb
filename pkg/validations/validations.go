package validations

import (
	"fmt"
	"regexp"
	"slices"
)

func ValidateDatabaseEngine(engine string) error {
	validEngines := []string{"postgres", "mysql", "mongodb", "cockroachdb", "sqlite"}
	valid := slices.Contains(validEngines, engine)
	if !valid {
		return fmt.Errorf("invalid database engine select from ('postgres', 'mysql', 'mongodb', 'cockroachdb', 'sqlite')")
	}
	return nil
}

func ValidateDatabaseHost(host string) error {
	match, err := regexp.MatchString("^(([a-zA-Z0-9-]{1,63}\\.)+[a-zA-Z]{2,}|localhost|(\\d{1,3}\\.){3}\\d{1,3})$", host)
	if err != nil {
		return err
	}
	if !match {
		return fmt.Errorf("it should be a valid host (e.g.: 127.0.0.1|localhost|192.168.0.1|example.com)")
	}
	return nil
}
