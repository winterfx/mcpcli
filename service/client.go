package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

const startClientTimeout = 15 * time.Second

func CreateClient(ctx context.Context, server *McpServer) (*client.StdioMCPClient, error) {
	var env []string
	for k, v := range server.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	c, err := client.NewStdioMCPClient(
		server.Command,
		env,
		server.Args...)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP client: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), startClientTimeout)
	defer cancel()
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "mcphost",
		Version: "0.1.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}
	fmt.Println("Initializing MCP client...")
	_, err = c.Initialize(ctx, initRequest)
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("failed to initialize MCP client: %w", err)
	}
	return c, nil
}

func RetrieveTools(ctx context.Context, c *client.StdioMCPClient) ([]mcp.Tool, error) {
	tools, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		err := parseMcpError(err)
		if errors.Is(err, ErrMethodNotFound) {
			return nil, fmt.Errorf("the server does not support tools: %w", err)
		} else {
			return nil, fmt.Errorf("failed to retrieve tools: %w", err)
		}
	}
	return tools.Tools, nil
}

func RetrievePrompts(ctx context.Context, c *client.StdioMCPClient) ([]mcp.Prompt, error) {
	p, err := c.ListPrompts(ctx, mcp.ListPromptsRequest{})
	if err != nil {
		err := parseMcpError(err)
		if errors.Is(err, ErrMethodNotFound) {
			return nil, fmt.Errorf("the server does not support prompts: %w", err)
		} else {
			return nil, fmt.Errorf("failed to retrieve prompts: %w", err)
		}
	}

	return p.Prompts, nil
}

func RetrieveResources(ctx context.Context, c *client.StdioMCPClient) ([]mcp.Resource, error) {
	r, err := c.ListResources(ctx, mcp.ListResourcesRequest{})
	if err != nil {
		err := parseMcpError(err)
		if errors.Is(err, ErrMethodNotFound) {
			return nil, fmt.Errorf("the server does not support resources: %w", err)
		} else {
			return nil, fmt.Errorf("failed to retrieve resources: %w", err)
		}
	}
	return r.Resources, nil
}

func CallTool(ctx context.Context, c *client.StdioMCPClient, toolName string, argsMap map[string]any) (*mcp.CallToolResult, error) {
	request := mcp.CallToolRequest{}
	request.Params.Name = toolName
	request.Params.Arguments = argsMap
	fmt.Printf("Calling tool '%s' with arguments: %v\n", toolName, argsMap)
	result, err := c.CallTool(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to call tool '%s': %v", toolName, err)
	}

	return result, nil

}
