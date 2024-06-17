package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

const (
	// PORT specifies the port where the server will be listening
	// Required for dev, and prod
	PORT = "PORT"

	// ENV specifies the environment where the server will be running
	// Required for dev, and prod
	ENV = "ENV"

	// DB_URI specifies the URI for the database
	// Required for dev, and prod
	DB_URI = "DB_URI"
)

var (
	// EnvVars is the struct holding the environment variables
	EnvVars *EnvConfig = newEnvConfig()

	// Required is the list of required environment variables for all environments
	required = []string{
		DB_URI,
	}

	// DevRequired is the list of required environment variables for the development environment
	devRequired = []string{}

	// ProdRequired is the list of required environment variables for the production environment
	prodRequired = []string{}
)

var (
	portDefault = "8080"
	envDefault  = "development"
)

// EnvConfig is the configuration for the environment variables
type EnvConfig struct {
	env map[string]string
}

// newEnvConfig creates a new EnvConfig
func newEnvConfig() *EnvConfig {
	env, err := godotenv.Read()
	if err != nil {
		fmt.Println("No .env file found.")
	}

	config := &EnvConfig{
		env: env,
	}
	config.setDefaults()
	config.validate()

	return config
}

// validate checks if the environment variables are set
func (e *EnvConfig) validate() {
	// Validate required environment variables
	e.validateRequired(required)
	environment := e.GetEnv()

	// Validate environment specific environment variables
	switch environment {
	case "development":
		e.validateRequired(devRequired)
	case "production":
		e.validateRequired(prodRequired)
	default:
		panic(fmt.Sprintf("Invalid environment: %s", environment))
	}
}

// validateRequired checks if the required environment variables are set
func (e *EnvConfig) validateRequired(required []string) {
	for _, key := range required {
		if e.env[key] == "" {
			panic(fmt.Sprintf("%s environment variable is not set", key))
		}
	}
}

// setDefaults sets the default values for the environment variables
func (e *EnvConfig) setDefaults() {
	if e.env[PORT] == "" {
		e.env[PORT] = portDefault
	}

	if e.env[ENV] == "" {
		e.env[ENV] = envDefault
	}
}

// GetPort returns the port where the server will be listening
func (e *EnvConfig) GetPort() string {
	return e.env[PORT]
}

// GetEnv returns the environment where the server will be running
func (e *EnvConfig) GetEnv() string {
	return e.env[ENV]
}

// GetDBURI returns the URI for the database
func (e *EnvConfig) GetDBURI() string {
	return e.env[DB_URI]
}

// Get returns the value of a specific environment variable
func (e *EnvConfig) Get(key string) string {
	if e.env[key] == "" {
		panic(fmt.Sprintf("%s environment variable is not set", key))
	}

	return e.env[key]
}
