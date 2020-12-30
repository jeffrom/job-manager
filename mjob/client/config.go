package client

// Config holds configuration data for the job client.
type Config struct {
	Host string `envconfig:"host" json:"host"`
}

var ConfigDefaults = Config{
	Host: ":1874",
}
