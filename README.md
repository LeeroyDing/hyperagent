# Hyperagent: The High-Agency OS Companion

![Hyperagent Header](https://img.shields.io/badge/Status-Production--Ready-green?style=for-the-badge) ![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)

**Hyperagent** is a next-generation, high-agency autonomous AI agent built in Go. Designed for power users and system administrators, Hyperagent moves beyond the sandbox, providing a direct, high-performance interface between the Gemini Pro LLM and your host operating system.

---

## 🚀 Core Philosophy

Hyperagent is built on the principle of **Direct Agency**. Unlike traditional agents that live in isolated containers, Hyperagent is a "Venom suit" for your machine. It runs as a native process, allowing it to manage your files, services, and workflows with zero latency and maximum permission depth.

## ✨ Key Features

### 1. Native Host Control
Hyperagent executes commands directly on your shell. It has full access to your local environment, CLI tools, and system configurations. It doesn't just suggest code; it maintains your system.

### 2. Model Context Protocol (MCP) Support
Built-in support for the **Model Context Protocol (MCP)** allows Hyperagent to instantly connect to a global ecosystem of tools. Whether it's interacting with Google Maps, Slack, GitHub, or local SQLite databases, Hyperagent consumes MCP servers as native capabilities.

### 3. Local-First Vector Memory
Using **chromem-go**, Hyperagent features a built-in, zero-dependency vector database. It generates embeddings using Gemini’s `text-embedding-004` and stores them locally, ensuring your long-term memory is private, fast, and persistent without needing external database clusters.

### 4. Gemini-Optimized Architecture
Hyperagent is engineered specifically for the **Google Gemini API**. It leverages Gemini's massive 1M+ token context window to maintain deep awareness of your system state, tool definitions, and conversation history simultaneously.

---

## 🏗️ Architecture Overview

Hyperagent is a single, statically linked Go binary. Its internal loop is optimized for low-latency reasoning and execution:

| Component | Implementation | Purpose |
| :--- | :--- | :--- |
| **The Brain** | Gemini 1.5 Pro / Flash | Reasoning, planning, and tool selection. |
| **The Memory** | `chromem-go` | Local vector storage for long-term recall. |
| **The Interface** | `mark3labs/mcp-go` | Standardized tool execution via MCP. |
| **The Hands** | `os/exec` | Direct shell and system command execution. |
| **The Storage** | JSONL / SQLite | Built-in session and state persistence. |

---

## 🛠️ Installation & Setup

### Prerequisites
*   **Go 1.25+**
*   **Gemini API Key** (from Google AI Studio)

### Quick Start
1.  **Clone and Build:**
    ```bash
    git clone https://github.com/LeeroyDing/hyperagent
    cd hyperagent
    go build -o hyperagent main.go
    ```

2.  **Configure Environment:**
    Create a `.env` file in the root directory:
    ```env
    GEMINI_API_KEY=your_api_key_here
    HYPERAGENT_LOG_LEVEL=info
    ```

3.  **Run:**
    ```bash
    ./hyperagent
    ```

---

## 🔧 Tooling & Extensibility

### Built-in Tools
Hyperagent ships with a core set of high-privilege tools:
*   **Shell**: Execute any command on the host system.
*   **Recall**: Query the local vector database for past interactions.
*   **Memorize**: Explicitly save information to long-term memory.

### Adding MCP Servers
To extend Hyperagent, simply add your MCP server configurations to `config.yaml`:
```yaml
mcp_servers:
  - name: github
    command: npx
    args: ["-y", "@modelcontextprotocol/server-github"]
```

---

## 🛡️ Security & Responsibility

**Hyperagent is a high-privilege tool.** 
By design, it does not run in a sandbox. It has the power to modify your system, delete files, and manage processes. 

*   **User Supervision**: It is recommended to run Hyperagent in `interactive` mode for destructive commands.
*   **Local-Only**: All memory and logs stay on your machine. No data is sent to third-party vector providers.

---

## 📈 Roadmap
*   [x] Core Go-Gemini Loop
*   [x] Local Vector Memory (chromem-go)
*   [x] MCP Client Integration
*   [ ] Web-based Dashboard (Localhost)
*   [ ] Advanced Multi-step Planning Visualizer

**Hyperagent** — *Your machine, amplified.*
