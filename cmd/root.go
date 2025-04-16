package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/winterfx/mcpcli/service"

	"github.com/spf13/cobra"
)

var (
	configFile string
	mcpConfig  *service.McpConfig
)

var rootCmd = &cobra.Command{
	Use:   "mcp-cli",
	Short: "A command line tool for interacting with MCP server",
	Long:  `A command line tool for managing MCP servers through an interactive interface`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use `mcp-cli [command] --help` for available commands.")
	},
}

func initConfig() {
	var err error
	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("Config file does not exist: %s\n", configFile)
		os.Exit(1)
	}
	mcpConfig, err = service.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config file: %v\n", err)
		os.Exit(1)
	}

}

func init() {
	cobra.OnInitialize(initConfig)
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	configFile = filepath.Join(home, ".mcp-cli.json")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", configFile, "config file (default is $HOME/.mcp-cli.json)")
	rootCmd.AddCommand(serverCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
