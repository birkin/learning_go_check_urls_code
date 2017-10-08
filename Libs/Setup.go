package Libs


func load_settings( settings Settings ) Settings {
	/* Loads settings, eventually for logging and database.
	   Called by main() */
	err := envconfig.Process("url_check_", &settings) // env settings look like `URL_CHECK__THE_SETTING`
	if err != nil {
		msg := fmt.Sprintf("error loading settings, ```%v```", err)
		rlog.Error(msg)
		panic(msg)
	}
	rlog.Debug(fmt.Sprintf("settings after settings initialized, ```%#v```", settings))
	return Settings{}
}
