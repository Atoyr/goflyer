package models

type config struct {
	apppath     string
	apikey      string
	timeoutmsec int64
}

// GetConfig is Getting config
// if path is empty this use default config path
func GetConfig(path string) config {
	var c config
	apppath := path
	if apppath == "" {
		apppath = ""
	}

	c.apppath = apppath
	return c
}
