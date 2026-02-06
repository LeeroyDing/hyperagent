package mcp

import (
"context"
"fmt"
"github.com/mark3labs/mcp-go/client"
"github.com/mark3labs/mcp-go/mcp"
)

type ServerConfig struct {
Name    string   `yaml:"name"` 
Command string   `yaml:"command"` 
Args    []string `yaml:"args"` 
Env     []string `yaml:"env"` 
}

type MCPManager struct {
Clients map[string]*client.Client
Tools   map[string]mcp.Tool
}

func NewMCPManager() *MCPManager {
return &MCPManager{
Clients: make(map[string]*client.Client),
Tools:   make(map[string]mcp.Tool),
}
}

func (m *MCPManager) AddServer(ctx context.Context, config ServerConfig) error {
c, err := client.NewStdioMCPClient(config.Command, config.Env, config.Args...)
if err != nil {
return fmt.Errorf("failed to create MCP client %s: %w", config.Name, err)
}

// Initialize the client
_, err = c.Initialize(ctx, mcp.InitializeRequest{
Params: mcp.InitializeParams{
ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
Capabilities:    mcp.ClientCapabilities{},
ClientInfo: mcp.Implementation{
Name:    "hyperagent",
Version: "0.1.0",
},
},
})
if err != nil {
return fmt.Errorf("failed to initialize MCP client %s: %w", config.Name, err)
}

m.Clients[config.Name] = c

// Discover tools
toolsResp, err := c.ListTools(ctx, mcp.ListToolsRequest{})
if err != nil {
return fmt.Errorf("failed to list tools for %s: %w", config.Name, err)
}

for _, tool := range toolsResp.Tools {
m.Tools[tool.Name] = tool
}

return nil
}

func (m *MCPManager) CallTool(ctx context.Context, serverName, toolName string, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
client, ok := m.Clients[serverName]
if !ok {
return nil, fmt.Errorf("MCP server %s not found", serverName)
}

return client.CallTool(ctx, mcp.CallToolRequest{
Params: mcp.CallToolParams{
Name:      toolName,
Arguments: arguments,
},
})
}
