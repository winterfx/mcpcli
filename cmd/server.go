package cmd

import (
	_ "bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/winterfx/mcpcli/service"

	"github.com/chzyer/readline"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage MCP servers",
	Long:  `Manage MCP servers`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use `mcp-cli server [command] --help` for available commands.")
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all servers",
	Long:  `List all servers`,
	Run: func(cmd *cobra.Command, args []string) {
		var names []string
		for name := range mcpConfig.MCPServers {
			names = append(names, name)
		}
		sort.Strings(names)

		serversView(mcpConfig, names)
	},
}

var inspectServerName string

// 添加这个辅助函数用于解析命令行
func parseCommandLine(input string) (cmd string, args []string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", nil
	}

	// 找到第一个空格，分离命令
	if idx := strings.Index(input, " "); idx != -1 {
		cmd = input[:idx]
		remaining := strings.TrimSpace(input[idx+1:])

		// 检查是否包含 JSON
		if strings.Contains(remaining, "{") {
			// 找到工具名称
			if spaceIdx := strings.Index(remaining, " "); spaceIdx != -1 {
				toolName := strings.TrimSpace(remaining[:spaceIdx])
				jsonPart := strings.TrimSpace(remaining[spaceIdx+1:])
				args = []string{toolName, jsonPart}
			} else {
				args = []string{remaining}
			}
		} else {
			// 普通参数模式
			args = strings.Fields(remaining)
		}
	} else {
		cmd = input
	}
	return cmd, args
}

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect a server",
	Long:  `Inspect a server`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		if inspectServerName == "" {
			fmt.Println("Error: server name is required")
			return
		}

		server, exists := mcpConfig.MCPServers[inspectServerName]
		if !exists {
			fmt.Printf("Error: server '%s' not found\n", inspectServerName)
			return
		}
		fmt.Printf("Inspecting server: %+v\n", server)
		client, err := service.CreateClient(ctx, &server)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			return
		}
		fmt.Printf("Entering interactive shell for server: %s\n", inspectServerName)
		fmt.Printf("Type 'help' for available commands, 'exit' to quit\n\n")

		rl, err := readline.NewEx(&readline.Config{
			Prompt:          fmt.Sprintf("%s> ", inspectServerName),
			HistoryFile:     "/tmp/mcp-cli.history",
			InterruptPrompt: "^C",
			EOFPrompt:       "exit",
		})
		if err != nil {
			fmt.Printf("Error initializing readline: %v\n", err)
			return
		}
		defer rl.Close()

		for {
			line, err := rl.Readline()
			if err != nil {
				if err == readline.ErrInterrupt {
					continue
				} else if err == io.EOF {
					break
				}
				fmt.Printf("Error reading input: %v\n", err)
				continue
			}

			cmd, args := parseCommandLine(line)
			if cmd == "" {
				continue
			}

			switch cmd {
			case "exit", "quit":
				return
			case "help":
				fmt.Println("Available commands:")
				fmt.Println("  tools    - Show server tools")
				fmt.Println("  prompts  - Show server status")
				fmt.Println("  resources  - Show server status")
				fmt.Println("  help    - Show this help")
				fmt.Println("  exit    - Exit the shell")
			case "info":
				fmt.Printf("Command: %s\n", server.Command)
				fmt.Printf("Arguments: %s\n", strings.Join(server.Args, " "))
				fmt.Println("\nEnvironment Variables:")
				for k, v := range server.Env {
					fmt.Printf("  %s = %s\n", k, v)
				}
			case "tools":
				tools, err := service.RetrieveTools(ctx, client)
				if err != nil {
					fmt.Printf("Error listing tools: %v\n", err)
					continue
				}
				toolsView(tools)
			case "prompts":
				prompts, err := service.RetrievePrompts(ctx, client)
				if err != nil {
					fmt.Printf("Error listing prompts: %v\n", err)
					continue
				}
				if len(prompts) == 0 {
					fmt.Println("No prompts found")
					continue
				}
				promptsView(prompts)
			case "resources":
				resources, err := service.RetrieveResources(ctx, client)
				if err != nil {
					fmt.Printf("Error listing resources: %v\n", err)
					continue
				}
				if len(resources) == 0 {
					fmt.Println("No resources found")
					continue
				}
				resourcesView(resources)
			case "call":
				if len(args) == 0 {
					fmt.Println("Usage: call <tool> '{...json...}'")
					continue
				}

				toolName := args[0]
				toolArgs := args[1:]
				argsMap := make(map[string]any)

				if len(toolArgs) == 1 && strings.HasPrefix(toolArgs[0], "{") {
					// JSON 参数模式
					if err := json.Unmarshal([]byte(toolArgs[0]), &argsMap); err != nil {
						fmt.Printf("Invalid JSON format: %v\n", err)
						continue
					}
				} else {
					fmt.Printf("Error: Invalid argument format. Usage: call <tool> '{\"param1\": \"value1\", ...}'\n")
					continue
				}

				result, err := service.CallTool(ctx, client, toolName, argsMap)
				if err != nil {
					fmt.Printf("Error calling tool '%s': %v\n", toolName, err)
					continue
				}

				fmt.Println("Result:")
				fmt.Println(result)

			default:
				fmt.Printf("Unknown command: %s\nType 'help' for available commands\n", cmd)
			}
		}
	},
}

func toolsView(tools []mcp.Tool) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"NAME", "DESCRIPTION", "PARAMETERS"})

	// Get terminal width
	width := 120 // default total width
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		width = w
	}

	// Calculate column widths based on terminal width
	nameWidth := width / 5
	descWidth := width / 5
	paramWidth := width / 3 // Parameters get more space
	fmt.Printf("total %d tools\n", len(tools))
	for _, tool := range tools {
		params, err := json.MarshalIndent(tool.InputSchema, "", "··")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		t.AppendRow(table.Row{
			tool.Name,
			tool.Description,
			string(params),
		})
	}

	// Configure table style
	t.SetStyle(table.StyleLight)
	t.Style().Options.DrawBorder = true
	t.Style().Options.SeparateColumns = true
	t.Style().Options.SeparateRows = true
	// Configure column constraints
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMax: nameWidth},
		{Number: 2, WidthMax: descWidth},
		{Number: 3, WidthMax: paramWidth},
	})

	t.Render()
}

func promptsView(prompts []mcp.Prompt) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"NAME", "DESCRIPTION"})
	for _, prompt := range prompts {
		t.AppendRow(table.Row{prompt.Name, prompt.Description})
	}
	t.SetStyle(table.StyleLight)
	t.Render()
}

func resourcesView(resources []mcp.Resource) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"NAME", "DESCRIPTION"})
	for _, resource := range resources {
		t.AppendRow(table.Row{resource.Name, resource.Description})
	}
	t.SetStyle(table.StyleLight)
	t.Render()
}

func serversView(mcpConfig *service.McpConfig, names []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Server Name", "Command", "Args", "Env"})

	for _, name := range names {
		s := mcpConfig.MCPServers[name]
		argStr := strings.Join(s.Args, " ")

		envParts := make([]string, 0, len(s.Env))
		for k, v := range s.Env {
			envParts = append(envParts, fmt.Sprintf("%s=%s", k, v))
		}
		sort.Strings(envParts)
		envStr := strings.Join(envParts, " ")
		t.AppendRow(table.Row{name, s.Command, argStr, envStr})
	}
	t.SetStyle(table.StyleLight)
	t.Render()
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(listCmd, inspectCmd)
	inspectCmd.Flags().StringVarP(&inspectServerName, "name", "n", "", "Name of the server to inspect")
}
