# hoho

<p>
  <img src="assets/logo.svg" width="540" alt="hoho — A modern Python debugger TUI">
</p>

A modern Python debugger TUI that replaces ipdb.

Built with Go, [Bubble Tea](https://github.com/charmbracelet/bubbletea), and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

## Prerequisites

- Python 3.8+
- [Go](https://golang.org/doc/install) toolchain (required for the first run as it compiles the Go TUI locally)

## Install

### pip / uv

```bash
uv pip install hoho
# or
pip install hoho
```

## Usage

Simply import `hoho` and set a trace wherever you want to pause execution and inspect state.

```python
import hoho

def calculate_metrics(data):
    total = sum(data)
    average = total / len(data)
    
    # Execution will pause here and the TUI will launch
    hoho.set_trace() 
    
    return average
```

Each invocation opens the Go-based TUI directly in your terminal, pausing Python execution. The interface presents:
- **Variables Pane**: See all local and global variables in scope.
- **Terminal Pane**: An interactive prompt where you can step through code.

### Keybindings

| Key | Action |
|---|---|
| `n` / `next` / `c` | Continue to next statement / Exit TUI |
| `q` / `Esc` / `Ctrl+C` | Quit |

## How it Works

The `hoho` package contains a lightweight JSON-RPC server. When `set_trace()` is hit, the Python process pauses, opens a local socket, and launches the compiled Go TUI as a subprocess. The Go TUI connects to the socket, receives the serialized state (including native DataFrame support), and renders the modern interface.
