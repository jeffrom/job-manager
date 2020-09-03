package jobclient

// Config holds configuration data for the job client.
type Config struct {
	Addr string `envconfig:"addr" json:"addr"`
}

var ConfigDefaults = Config{
	Addr: ":1874",
}
