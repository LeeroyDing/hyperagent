# Hyperagent Alpha Release Notes

## Version 0.1.0-alpha

We are excited to announce the Alpha release of Hyperagent! This release focuses on stability, core functionality, and engineering excellence.

### Key Features
- **Autonomous Agentic Loop**: Powered by Gemini 1.5 Pro/Flash.
- **Direct Host Control**: Execute shell commands safely and effectively.
- **Local Vector Memory**: Long-term memory using chromem-go.
- **MCP Support**: Extensible tool usage via Model Context Protocol.

### Changes & Improvements
- **Streamlined Scope**: Removed TUI, WASM, Telegram, and PKM Mirroring to focus on core stability.
- **Engineering Excellence**: Added golangci-lint and Makefile for better development workflows.
- **Improved Testing**: Added unit tests and integration tests.

### Installation
Clone the repository and run:
bash
make build


### Known Issues
- **Port Conflict**: The Web UI runs on port 3001, which may conflict with other services like the Gemini Web UI backend. Please ensure port 3001 is free before running.
