package backend

type Config struct {
	// MaxStreamSize is the maximum number of stream events to store.
	MaxStreamSize int

	// HistoryLimit is the maximum number of resource versions to store.
	HistoryLimit int
}

var DefaultConfig = Config{
	MaxStreamSize: 100000,
	HistoryLimit:  10,
}
