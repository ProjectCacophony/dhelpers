package dhelpers

import (
	"os"
	"strings"
)

// Environment is the type for possible Environments
type Environment string

const (
	// EnvironmentDevelopment is the development environment
	EnvironmentDevelopment Environment = "development"
	// EnvironmentTesting is the testing environment
	EnvironmentTesting Environment = "testing"
	// EnvironmentStaging is the staging environment
	EnvironmentStaging Environment = "staging"
	// EnvironmentProduction is the production environment
	EnvironmentProduction Environment = "production"
)

// GetEnvironment returns the current environment
// it reads the environment from the Environment variable ENV or Environment
// if none is set it is always development
func GetEnvironment() Environment {
	switch strings.ToLower(os.Getenv("ENV")) {
	case "testing", "test":
		return EnvironmentTesting
	case "staging", "stag":
		return EnvironmentStaging
	case "production", "prod":
		return EnvironmentProduction
	}

	switch strings.ToLower(os.Getenv("Environment")) {
	case "testing", "test":
		return EnvironmentTesting
	case "staging", "stag":
		return EnvironmentStaging
	case "production", "prod":
		return EnvironmentProduction
	}

	return EnvironmentDevelopment
}
