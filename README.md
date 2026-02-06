# Hyperagent

Hyperagent is a high-agency, Go-native autonomous AI agent designed for direct host system control. It leverages the Gemini 1.5 API for reasoning and the Model Context Protocol (MCP) for tool extensibility.

## Features
- **Core Agentic Loop**: Autonomous planning and execution using Gemini 1.5 Pro/Flash.
- **Direct Shell Execution**: Control your host system directly via shell commands.
- **Local Vector Memory**: Persistent long-term memory using `chromem-go` and Gemini embeddings.
- **MCP Support**: Seamlessly integrate and use tools from any MCP-compatible server.
- **Session Persistence**: Conversation history is automatically saved to `history.jsonl`.
- **Robust Error Handling**: Improved parsing and tool execution resilience.

## Installation
Ensure you have Go 1.25+ installed.

```bash
go build -o hyperagent main.go
```

## Usage
Set your Gemini API key:

```bash
export GEMINI_API_KEY=your_api_key_here
```

Run the agent with a prompt:

```bash
./hyperagent "Your request here"
```

## Configuration
Configure MCP servers in `config.yaml`:

```yaml
mcp_servers:
  - name: filesystem
    command: npx
    args: ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/allowed/dir"]
```

## License
MIT
