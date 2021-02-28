package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Environment struct {
	Port              string
}

func Get() (*Environment, error) {
	var env Environment
	var missed []string

	for _, tmp := range []struct {
		field *string
		name  string
	}{
		{&env.Port, "PORT"},
	} {
		v := os.Getenv(tmp.name)
		if v == "" {
			missed = append(missed, tmp.name)
		}
		fmt.Println(v)
		*tmp.field = v
	}

	if 0 < len(missed) {
		return nil, errors.New(fmt.Sprintf("%s cannot be empty.", strings.Join(missed, ", ")))
	}

	return &env, nil
}
