package service

import (
	"encoding/json"
	"fmt"
	"os"
)

type McpConfig struct {
	Path       string               `json:"-"`
	MCPServers map[string]McpServer `json:"mcpServers"`
}

type McpServer struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env,omitempty"`
}

func LoadConfig(path string) (*McpConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config McpConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	config.Path = path
	return &config, nil
}

// SaveConfig saves the current configuration to the specified file
func (c *McpConfig) SaveConfig() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(c.Path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
