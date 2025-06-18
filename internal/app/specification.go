package app

// Specification holds the configuration for the application.
type Specification struct {
	Address     string `envconfig:"address" default:":8080"`
	DatabaseURL string `envconfig:"database_url" required:"true"`
}
