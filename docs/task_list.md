# Hyperagent Implementation Task List

This document provides a granular breakdown of tasks for each phase, including verification mechanisms.

## Phase 1: Core Infrastructure & Lifecycle

### Tasks
1.  **CLI Framework**: Integrate spf13/cobra in main.go to replace standard flags with a subcommand structure (up, down, status, version).
2.  **Daemon Package**: Create internal/daemon to handle:
    - PID file management in ~/.hyperagent/hyperagent.pid.
    - Signal handling (SIGTERM/SIGINT) for graceful shutdown.
    - Backgrounding logic (forking or using a library like sevlyar/go-daemon).
3.  **Daemon API**: Refactor internal/web/server.go to provide a "Daemon API" on a configurable port (default 8080) for CLI-to-Daemon communication.
4.  **Command Implementation**:
    - up: Check for existing PID, start daemon, initialize API.
    - down: Read PID, send SIGTERM, wait for process exit, cleanup PID.
    - status: Ping the Daemon API and check PID file health.

### Verification
- **Automated**: 
  - Unit tests in internal/daemon for PID locking and stale PID cleanup.
  - Integration test script: hyperagent up && hyperagent status && hyperagent down verifying exit codes and process existence.
- **Manual**:
  - Run hyperagent up, verify process with ps aux | grep hyperagent.
  - Verify ~/.hyperagent/hyperagent.pid contains the correct PID.

## Phase 2: Agent Management & Continuous Persistence

### Tasks
1.  **Agent Registry**: Implement a thread-safe registry in the daemon to track active agent.Agent instances by ID.
2.  **Continuous Persistence Hook**: Modify internal/agent/agent.go (Run and handleToolCall methods) to trigger a state save to ~/.hyperagent/state/<session_id>.json after every tool execution and model response.
3.  **State Schema**: Define a JSON schema for agent state including conversation history, current tool context, and metadata.
4.  **Management Commands**:
    - start: CLI sends request to Daemon API -> Daemon creates Agent -> Returns ID.
    - list: CLI queries Daemon API for registry contents.
    - rm: CLI sends delete request -> Daemon stops agent and deletes state files.
5.  **Recovery Logic**: Update daemon startup to scan ~/.hyperagent/state/ and re-instantiate agents found in the directory.

### Verification
- **Automated**:
  - Test in internal/agent verifying that Run() results in a file write to the state directory.
  - Mock recovery test: Manually place a state file and verify the daemon loads it into the registry on startup.
- **Manual**:
  - Start an agent, perform one action, kill the daemon with kill -9, run hyperagent up, and verify the agent is still in hyperagent list with its history intact.

## Phase 3: Observability & Interaction

### Tasks
1.  **Log Redirection**: Update the daemon to capture stdout/stderr of each agent and pipe it to ~/.hyperagent/logs/<session_id>.log.
2.  **Log Tailer**: Implement hyperagent logs -f <id> using a tailing library to stream the agent's log file.
3.  **Interactive Bridge**: Implement a WebSocket or Unix Socket bridge in the Daemon API to allow the attach command to send input to a background agent's interactive confirmation prompt.
4.  **Doctor Command**: Implement internal/doctor to check:
    - Gemini API connectivity.
    - Disk write permissions for state/logs.
    - Presence of required binaries (e.g., shell).

### Verification
- **Automated**:
  - Test log file rotation and creation.
  - doctor unit tests with mocked failures (e.g., missing API key).
- **Manual**:
  - Run an agent in the background, run hyperagent logs -f <id>, and watch the 'Think-Act' loop in real-time.

## Phase 4: Automation & Portability

### Tasks
1.  **Cron Integration**: Add robfig/cron to the daemon to handle scheduled tasks.
2.  **Schedule Command**: Implement CLI to register cron expressions and commands with the daemon's cron runner.
3.  **Export/Import**: Implement tar.gz bundling of an agent's state directory (history, memory, state.json).
4.  **Prune Logic**: Implement a background worker in the daemon to cleanup agents marked as 'stopped' or older than a configurable retention period.

### Verification
- **Automated**:
  - Test export/import by exporting an agent, deleting it, importing it, and verifying hash of the state file.
  - Cron trigger test with a 1-minute interval.
- **Manual**:
  - Schedule a task for 2 minutes in the future, verify it starts and completes via hyperagent list.

## Phase 5: Advanced Monitoring & Evals

### Tasks
1.  **Resource Monitor**: Use shirou/gopsutil to track CPU and Memory usage of the daemon and its sub-processes.
2.  **Top Command**: Implement a TUI (Terminal User Interface) for hyperagent top showing live agent stats.
3.  **Eval Runner**: Create a framework to run predefined prompts and compare agent outputs against expected 'golden' responses.

### Verification
- **Automated**:
  - Unit tests for resource metric collection.
- **Manual**:
  - Run hyperagent top and verify it reflects active agent resource consumption.
