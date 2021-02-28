package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Environment struct {
	Port                        string
	LineChannelSecret           string
	LineChannelToken            string
	FirebaseCredentialsFilePath string
	FirebaseDatabaseURL         string
}

func Get() (*Environment, error) {
	var env Environment
	var missed []string

	for _, tmp := range []struct {
		field *string
		name  string
	}{
		{&env.Port, "PORT"},
		{&env.LineChannelSecret, "LINE_CHANNEL_SECRET"},
		{&env.LineChannelToken, "LINE_CHANNEL_TOKEN"},
		{&env.FirebaseCredentialsFilePath, "FIREBASE_CREDENTIALS_FILE_PATH"},
		{&env.FirebaseDatabaseURL, "FIREBASE_DATABASE_URL"},
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
