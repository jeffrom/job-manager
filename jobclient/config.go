package jobclient

// Config holds configuration data for the job client.
type Config struct {
	Addr string `json:"addr"`
}

func getDefaultConfig() Config {
	return Config{}
}
