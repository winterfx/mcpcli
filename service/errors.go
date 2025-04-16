package service

import (
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

var ErrMethodNotFound = errors.New("mcp: method not found")

type mcpWrappedError struct {
	code  int
	msg   string
	inner error
}

func (e *mcpWrappedError) Error() string {
	return fmt.Sprintf("MCP error (%d): %s", e.code, e.msg)
}

func (e *mcpWrappedError) Unwrap() error {
	// 对应 JSON-RPC 标准错误码
	switch e.code {
	case mcp.METHOD_NOT_FOUND: // Method not found
		return ErrMethodNotFound
	// 你可以加更多的case，比如 -32700, -32600 等
	default:
		return e.inner
	}
}

func parseMcpError(err error) error {
	raw := err.Error()
	switch raw {
	case "Method not found":
		return &mcpWrappedError{code: mcp.METHOD_NOT_FOUND, msg: raw, inner: err}
	case "Invalid request":
		return &mcpWrappedError{code: mcp.INVALID_REQUEST, msg: raw, inner: err}
	default:
		return fmt.Errorf("unrecognized MCP error: %w", err)
	}
}
