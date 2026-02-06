# Persistent Shell Sessions (PTY/TTY)

## Overview
Hyperagent is moving from a stateless command execution model to a stateful, session-based shell model. Each conversation (session) will have its own dedicated persistent shell process (e.g., `bash` or `sh`) running in a Pseudo-Terminal (PTY).

## Benefits
1. **Stateful Navigation**: The working directory (`cd`) persists across multiple turns in a conversation.
2. **Environment Persistence**: Environment variables (`export`), aliases, and sourced scripts stay active for the duration of the session.
3. **Interactive Tool Support**: Using a PTY allows Hyperagent to interact with tools that require a terminal (e.g., `git commit` prompts, `npm init`, or interactive CLI wizards).
4. **Background Processes**: Long-running tasks (like `go run main.go` for a web server) can be started and monitored across different messages.
5. **Terminal Fidelity**: Tools that detect TTY support (like `ls` with colors or `gh` CLI) will behave as if they are running in a real terminal.

## Architecture

### 1. Session Manager
- Maintains a registry (map) of `SessionID` to `ShellSession` objects.
- Handles lifecycle management: creating, retrieving, and killing shell processes.
- Ensures cleanup of orphaned processes on application exit.

### 2. Shell Session (PTY)
- Uses `github.com/creack/pty` to spawn a shell process with a master/slave PTY pair.
- **Input**: Commands are written directly to the PTY master's `stdin`.
- **Output**: A background goroutine reads from the PTY master and buffers the output.
- **Prompt Detection**: Uses a unique delimiter or shell prompt detection to signal when a command has finished executing.

### 3. Integration
- The `ShellExecutor` tool is updated to accept a `SessionID`.
- If no session exists for the ID, a new one is spawned.
- Commands are routed to the corresponding persistent shell.

## Security Considerations
- **Isolation**: While sessions share the same host user, they are logically separated by process. 
- **Resource Management**: Limits on the number of concurrent sessions and process timeouts should be implemented to prevent resource exhaustion.
- **Cleanup**: Automatic termination of idle sessions after a configurable timeout.
