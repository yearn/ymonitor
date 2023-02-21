package config

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"strings"
)

type Host struct {
	Name    string `json:"name"`
	Url     Url
	Network string         `json:"network"`
	Auth    Authentication `json:"auth,omitempty"`
}

type Authentication struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type Url struct {
	Url url.URL `json:"url"`
}

func (u *Url) UnmarshalJSON(data []byte) error {
	var v string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	url, err := url.Parse(extractUrl(string(v)))
	u.Url = *url
	return nil
}

func extractUrl(url string) string {
	// urls can be either specified directly or via environment variable names:
	// "$FOO_BAR" will load the contents of $FOO_BAR into the url
	splitted := strings.Split(url, "$")
	if len(splitted) == 1 {
		return url
	}
	envName := splitted[1]
	if value, ok := os.LookupEnv(envName); ok {
		return value
	}
	panic(errors.New("Couldn't find env variable " + envName))
}
