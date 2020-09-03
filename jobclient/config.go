package jobclient

// Config holds configuration data for the job client.
type Config struct {
	Addr string `json:"addr"`
}

var ConfigDefaults = Config{
	Addr: ":1874",
}

func getDefaultConfig() Config {
	return Config{}
}
