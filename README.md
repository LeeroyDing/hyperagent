# Hyperagent (Alpha)

Hyperagent is a high-agency, Go-native autonomous AI agent designed for direct host system control. It leverages the **Gemini 3 Flash** API for reasoning and the Model Context Protocol (MCP) for tool extensibility.

## Features
- **Core Agentic Loop**: Autonomous planning and execution using Gemini 3 Flash.
- **Web UI**: A modern, Tailwind-based Web UI for interaction (Port 3001).
- **Direct Shell Execution**: Control your host system directly via shell commands.
- **Local Vector Memory**: Persistent long-term memory using chromem-go and Gemini embeddings.
- **MCP Support**: Seamlessly integrate and use tools from any MCP-compatible server.
- **Session Persistence**: Conversation history is automatically saved to history.jsonl.
- **Robust Error Handling**: Improved parsing and tool execution resilience.
- **Engineering Excellence**: Includes golangci-lint integration and a Makefile for automated workflows.

## Installation
Ensure you have Go 1.25+ installed.

bash
make build


## Usage
Set your Gemini API key:

bash
export GEMINI_API_KEY=your_api_key_here


Run the agent with a prompt:

bash
./hyperagent "Your request here"


## Development

### Linting
bash
make lint


### Testing
bash
make test


## Configuration
Configure MCP servers in config.yaml:

yaml
mcp_servers:
  - name: filesystem
    command: npx
    args: ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/allowed/dir"]


## License
MIT
