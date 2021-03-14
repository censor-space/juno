package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Environment struct {
	Port                        string
	LineChannelSecret           string
	LineChannelToken            string
	FirebaseCredentialsFilePath string
	FirebaseDatabaseURL         string
    ThinkingTimeSec int
    CORSAllowOrigins []string
}

func Get() (*Environment, error) {
	var env Environment
	var missed []string
    var thinkingTime string
    var corsAllowOrigins string

	for _, tmp := range []struct {
		field *string
		name  string
	}{
		{&env.Port, "PORT"},
		{&env.LineChannelSecret, "LINE_CHANNEL_SECRET"},
		{&env.LineChannelToken, "LINE_CHANNEL_TOKEN"},
		{&env.FirebaseCredentialsFilePath, "FIREBASE_CREDENTIALS_FILE_PATH"},
		{&env.FirebaseDatabaseURL, "FIREBASE_DATABASE_URL"},
        {&thinkingTime, "THINKING_TIME_SEC"},
        {&corsAllowOrigins, "CORS_ALLOW_ORIGINS"},
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

    tts, err := strconv.Atoi(thinkingTime)
    if err != nil {
        return nil, err
    }

    corsaos := strings.Split(corsAllowOrigins, ",")
    for _, corsao := range corsaos {
        env.CORSAllowOrigins = append(env.CORSAllowOrigins, corsao)
    }
    env.ThinkingTimeSec = tts

	return &env, nil
}
