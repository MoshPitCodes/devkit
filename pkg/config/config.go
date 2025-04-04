package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Resources defines resource limits for a container.
type Resources struct {
	Memory string `mapstructure:"memory" yaml:"memory"` // e.g., "512m", "1g"
	CPU    string `mapstructure:"cpu" yaml:"cpu"`       // e.g., "0.5", "1"
}

// ContainerTemplate defines the structure for a reusable container definition.
type ContainerTemplate struct {
	Image         string            `mapstructure:"image" yaml:"image"`
	Ports         []string          `mapstructure:"ports" yaml:"ports"`             // Format: ["host:container", ...]
	Volumes       []string          `mapstructure:"volumes" yaml:"volumes"`         // Format: ["host:container", ...]
	Environment   map[string]string `mapstructure:"environment" yaml:"environment"` // Key-value pairs
	Command       []string          `mapstructure:"command" yaml:"command"`         // Command to run
	NetworkMode   string            `mapstructure:"networkMode" yaml:"networkMode,omitempty"`
	Resources     *Resources        `mapstructure:"resources" yaml:"resources,omitempty"`         // Optional resource limits
	RestartPolicy string            `mapstructure:"restartPolicy" yaml:"restartPolicy,omitempty"` // e.g., "no", "always", "on-failure"
	// Add more fields as needed (e.g., healthcheck)
}

// --- New Config Structure ---

// AgentConfig holds configuration specific to the devkit agent.
type AgentConfig struct {
	SocketPath string `mapstructure:"socketPath"`
	// LogLevel string `mapstructure:"logLevel"`
}

// CliConfig holds configuration specific to the devkit CLI.
type CliConfig struct {
	DefaultTimeoutSeconds int `mapstructure:"defaultTimeoutSeconds"`
	// LogLevel string `mapstructure:"logLevel"`
}

// Config holds the application configuration (excluding container templates).
type Config struct {
	Agent             AgentConfig `mapstructure:"agent"`
	Cli               CliConfig   `mapstructure:"cli"`
	TemplateDirectory string      `mapstructure:"templateDirectory"` // Optional override for template dir
}

// Cfg holds the global loaded configuration.
var Cfg *Config

const (
	defaultTemplateDir    = ".devkit/templates"
	defaultAgentSocket    = "/tmp/devkitd.sock"
	defaultCliTimeoutSecs = 60 // Increase default CLI command timeout
)

// LoadConfig reads devkit-specific configuration from files and environment variables.
func LoadConfig() (*Config, error) {
	v := viper.New()

	// --- Set Defaults ---
	v.SetDefault("agent.socketPath", defaultAgentSocket)
	v.SetDefault("cli.defaultTimeoutSeconds", defaultCliTimeoutSecs)
	v.SetDefault("templateDirectory", defaultTemplateDir) // Default template location

	// --- Configuration file paths ---
	v.AddConfigPath(".devkit") // Project-specific config
	home, err := os.UserHomeDir()
	if err == nil {
		v.AddConfigPath(filepath.Join(home, ".config", "devkit")) // User-specific config
	}
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Attempt to read the config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found; defaults/env vars will be used
	}

	// --- Environment variable handling ---
	v.SetEnvPrefix("DEVKIT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// --- Unmarshal the config into the Config struct ---
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// --- Post-processing/Validation (Optional) ---
	// Ensure the agent's template directory exists
	templateDir := config.TemplateDirectory
	if templateDir == "" {
		templateDir = defaultTemplateDir
	}
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		// Log a warning, but maybe don't fail? Or should we fail?
		// For now, let's return an error if we can't ensure the dir exists.
		return nil, fmt.Errorf("failed to create template directory %s: %w", templateDir, err)
	}
	// Ensure socket path directory exists?

	Cfg = &config
	return Cfg, nil
}

// --- New Template Loading Logic ---

// LoadContainerTemplate loads a single container template definition from its file.
func LoadContainerTemplate(name string) (*ContainerTemplate, error) {
	if Cfg == nil {
		// Should not happen if LoadConfig is called first, but defensive check.
		return nil, fmt.Errorf("devkit configuration not loaded")
	}

	templateDir := Cfg.TemplateDirectory
	if templateDir == "" {
		templateDir = defaultTemplateDir // Use default if not set in config
	}
	// Adjust template path to look inside the 'containers' subdirectory
	containerTemplateDir := filepath.Join(templateDir, "containers")
	templateFilePath := filepath.Join(containerTemplateDir, name+".yaml")

	data, err := os.ReadFile(templateFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("template '%s' not found at %s: %w", name, templateFilePath, err)
		}
		return nil, fmt.Errorf("failed to read template file %s: %w", templateFilePath, err)
	}

	var template ContainerTemplate
	if err := yaml.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("failed to parse template file %s: %w", templateFilePath, err)
	}

	// --- Template Validation (Optional) ---
	if template.Image == "" {
		return nil, fmt.Errorf("template '%s' is missing required 'image' field", name)
	}

	return &template, nil
}
