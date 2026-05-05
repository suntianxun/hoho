# Hoho Debugger Design

## Overview
A Python debugging tool (`hoho`) that provides a vivid, modern Text User Interface (TUI) built in Go using the Charm ecosystem. The tool seamlessly wraps Python execution via a `hoho.set_trace()` import, taking over the terminal exactly like `ipdb` but providing a much richer, multi-pane visual experience.

## Architecture

### 1. The Python Server (`hoho` package)
- Provides `hoho.set_trace()`.
- When invoked, the Python process binds to an ephemeral TCP port and acts as a custom JSON-RPC server.
- The Python process then launches the Go TUI executable (`hoho-tui`) as a subprocess, passing the TCP port as an environment variable or argument.
- The Python execution thread blocks (pauses), waiting for debugging commands via the TCP socket.
- Handles safe serialization of `locals()` and `globals()`.
- Natively supports serializing pandas/polars DataFrames into paginated/tabular JSON structures.

### 2. The Go TUI Client (`hoho-tui`)
- Built entirely in Go utilizing Charm (`bubbletea`, `lipgloss`, `bubbles`).
- Connects to the Python server via local TCP socket.
- Renders the UI and takes over terminal input.
- Translates user UI interactions into JSON commands sent to the Python server.
- When the user issues a continue command (e.g., `c` or `next`), the Go process exits, returning control to the standard Python stdout/stdin.

## User Interface (TUI) Layout
The UI is divided into 3 primary panes, styled with vivid `lipgloss` colors:

1. **Variables Pane (Left):**
   - A tree or list view showing all variables in the current scope.
   - Public and private variables are visually distinct (e.g., private variables dimmed or colored differently).
   - Displays variable names and truncated values.

2. **Data Viewer Pane (Top Right):**
   - Dedicated viewer for complex data structures like DataFrames.
   - Utilizes `charmbracelet/bubbles/table` for a spreadsheet-like view.
   - Automatically populates when the user selects a DataFrame from the Variables Pane.

3. **Terminal Pane (Bottom Right):**
   - An interactive prompt (using `bubbles/textinput` or similar) replicating the `ipdb` experience.
   - Users can type standard debug commands (`n`, `c`, `s`) or arbitrary Python code to be evaluated in the current context.

## Data Flow Protocol
The communication uses a custom, lightweight JSON-over-TCP protocol.

**1. Initialization:**
- Python: `{"type": "state", "vars": [{"name": "my_df", "type": "DataFrame", "is_private": false}, ...]}`

**2. Evaluation (e.g., viewing a DataFrame):**
- Go: `{"type": "eval", "expr": "my_df.head(50)"}`
- Python: `{"type": "eval_result", "data": {"columns": [...], "rows": [...]}}`

**3. Execution Control:**
- Go: `{"type": "command", "cmd": "next"}`
- Python processes the `next` command. The Python server unpauses, the TCP connection drops, and the Go TUI gracefully exits.

## Error Handling
- If the Go TUI binary is missing, `hoho.set_trace()` gracefully degrades to standard `pdb`.
- Network serialization errors in Python are caught and sent to the TUI as error strings to display in the Terminal pane.