# Hyperagent CLI Reference

This document describes the command-line interface for the Hyperagent daemon and agent management.

## Daemon Commands

### hyperagent up
Starts the Hyperagent daemon in the background.

```text
Usage: hyperagent up [flags]

Start the Hyperagent daemon process. This will initialize the background service
responsible for managing agents, tasks, and persistent sessions.

Persistence Note: On startup, the daemon will automatically scan the local 
storage (default ~/.hyperagent/state) and restore all agent instances to 
their last recorded state. Because Hyperagent uses continuous persistence, 
this includes recovery from unexpected shutdowns or crashes.

Flags:
  -d, --detach          Run daemon in background (default true)
      --port int        Port for the daemon API (default 8080)
      --config string   Path to config file (default "~/.hyperagent/config.yaml")
  -h, --help            help for up
```

### hyperagent down
Stops the Hyperagent daemon and all running agents.

```text
Usage: hyperagent down [flags]

Gracefully shut down the Hyperagent daemon. This will stop the background process.
Since agent states are persisted continuously during operation, this command 
simply ensures a clean exit and stops all active loops.

Flags:
      --force           Force immediate shutdown
  -h, --help            help for down
```

### hyperagent status
Shows the current status of the daemon.

```text
Usage: hyperagent status

Check if the Hyperagent daemon is running and display system health,
including API endpoint, uptime, and resource usage.

Flags:
  -h, --help            help for status
```

### hyperagent doctor
Check the local environment for common issues.

```text
Usage: hyperagent doctor

Check the local environment for common issues: 
- Validates GEMINI_API_KEY
- Checks connectivity to Gemini API
- Verifies vector database (chromem-go) health
- Checks for required OS permissions
```

### hyperagent version
Show version information.

```text
Usage: hyperagent version

Show the Hyperagent version, build date, and git commit hash.
```

## Agent Management Commands

### hyperagent list
Lists all managed agents/tasks.

```text
Usage: hyperagent list [flags]

Display a table of all agents currently managed by the daemon, including
their IDs, status, and creation time.

Flags:
  -a, --all             Show all agents (including stopped/failed)
  -q, --quiet           Only display agent IDs
      --format string   Output format: table, json, yaml (default "table")
  -h, --help            help for list
```

### hyperagent start <args>
Starts a new agent instance.

```text
Usage: hyperagent start [prompt] [flags]

Create and start a new agent instance with the specified prompt or configuration.
If no prompt is provided, it will start an interactive session.

Arguments:
  prompt                The initial instruction or goal for the agent

Flags:
  -n, --name string     Assign a custom name to the agent
  -p, --profile string  Use a specific agent profile (e.g., 'developer', 'researcher')
      --env stringArray Set environment variables (e.g., KEY=VALUE)
  -h, --help            help for start
```

### hyperagent stop <id>
Stops a specific agent by its ID or name.

```text
Usage: hyperagent stop [ID|NAME]

Stop a specific running agent instance. The agent's state remains persisted 
in the local state store for later resumption.

Arguments:
  ID|NAME               The unique identifier or name of the agent to stop

Flags:
  -h, --help            help for stop
```

### hyperagent rm <id>
Permanently remove an agent instance.

```text
Usage: hyperagent rm [ID|NAME] [flags]

Permanently remove one or more agent instances and their associated data, 
including conversation history and local state files.

Arguments:
  ID|NAME               The unique identifier or name of the agent to remove

Flags:
  -f, --force          Stop the agent if it is running before removing
  -h, --help           help for rm
```

### hyperagent schedule
Schedule an agent to start and perform a task at recurring intervals.

```text
Usage: hyperagent schedule "cron-expression" -- "command" [flags]

Schedule an agent to start and perform a task at recurring intervals.

Example:
  hyperagent schedule "0 9 * * 1-5" -- "Summarize my unread emails"

Flags:
  -n, --name string     Assign a name to the scheduled task
  -p, --profile string  Use a specific agent profile
  -h, --help            help for schedule
```

### hyperagent logs <id>
Fetch the logs of a specific agent.

```text
Usage: hyperagent logs [ID|NAME] [flags]

Fetch the logs of a specific agent. Useful for seeing the 'Think-Act' loop 
output without attaching to the session.

Flags:
  -f, --follow         Follow log output
      --tail int       Number of lines to show from the end (default 20)
```

### hyperagent attach <id>
Attach terminal to a running agent.

```text
Usage: hyperagent attach [ID|NAME]

Attach your terminal to a running agent's interactive session. This allows 
you to provide manual feedback or take over the loop.
```

## System & Maintenance

### hyperagent export
Bundle an agent's entire state into a compressed file.

```text
Usage: hyperagent export [ID|NAME] [output-file]

Bundle an agent's entire state into a compressed file. This includes its 
conversation history and vector memory snapshots.
```

### hyperagent import
Restore an agent from a previously exported bundle.

```text
Usage: hyperagent import [input-file]

Restore an agent from a previously exported bundle.
```

### hyperagent prune
Cleanup stopped agents and temporary files.

```text
Usage: hyperagent prune [flags]

Remove all stopped agents, temporary files, and old session histories 
that are no longer needed.

Flags:
  -f, --force          Do not prompt for confirmation
      --older-than     Only prune agents older than duration (e.g. 24h)
```

### hyperagent config
Manage settings.

```text
Usage: hyperagent config [subcommand]

View or update Hyperagent configuration settings.

Subcommands:
  get <key>          Display a specific config value
  set <key> <value>  Update a config value
  list               Show all current configuration
```

### hyperagent memory
Interact with vector memory.

```text
Usage: hyperagent memory [subcommand]

Interact directly with the agent's long-term vector memory.

Subcommands:
  search "query"     Search across all agent memories
  export <id>        Export an agent's memory to a file
  clear              Wipe all stored memories
```
