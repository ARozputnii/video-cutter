package app

// Config holds all configurable values
type Config struct {
	Port        int
	UploadDir   string
	FrontendDir string
}

// NewConfig returns the default configuration.
func NewConfig() Config {
	return Config{
		Port:        3000,
		UploadDir:   "uploads",
		FrontendDir: "internal/frontend",
	}
}
