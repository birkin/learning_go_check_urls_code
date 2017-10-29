package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/romana/rlog"
)

type Settings struct {
	DB_USERNAME       string `envconfig:"DB_USERNAME" required:"true"`
	DB_PASSWORD       string `envconfig:"DB_PASSWORD" required:"true"`
	DB_HOST           string `envconfig:"DB_HOST" required:"true"`
	DB_PORT           string `envconfig:"DB_PORT" required:"true"`
	DB_NAME           string `envconfig:"DB_NAME" required:"true"`
	TEST_EMAIL_STRING string `envconfig:"TEST_EMAIL_STRING" required:"true"`
	MAIL_HOST         string `envconfig:"TEST_MAIL_HOST" required:"true"`
}

func load_settings() Settings {
	/* Loads settings, currently for logging and database.
	   Called by main() */
	var settings Settings
	rlog.Debug(fmt.Sprintf("settings before calling envconfig, `%#v`", settings))
	err := envconfig.Process("url_check_", &settings) // env settings look like `URL_CHECK__THE_SETTING`
	if err != nil {
		msg := fmt.Sprintf("error loading settings, ```%v```", err)
		// rlog.Error(msg)
		panic(msg)
	}
	rlog.Debug(fmt.Sprintf("settings after calling envconfig, `%#v`", settings))
	return settings
}
