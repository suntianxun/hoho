# Hoho Debugger Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a modern Python debugger TUI using Go/Charm that replaces ipdb, featuring a variables pane and terminal.

**Architecture:** Python `set_trace` starts a TCP server, pauses execution, and spawns a Go subprocess. The Go binary connects to the TCP socket, runs a Charm-based TUI, and communicates via custom JSON protocol.

**Tech Stack:** Python 3, Go, Bubble Tea, Lip Gloss.

---

### Task 1: Python Package Setup and State Serialization

**Files:**
- Create: `hoho/__init__.py`
- Create: `hoho/serializer.py`
- Create: `tests/test_serializer.py`

- [ ] **Step 1: Write the failing test for serialization**

```python
# tests/test_serializer.py
from hoho.serializer import serialize_state

def test_serialize_locals():
    x = 1
    _y = "private"
    state = serialize_state(locals())
    
    assert len(state) >= 2
    vars_dict = {v["name"]: v for v in state}
    assert vars_dict["x"]["value"] == "1"
    assert vars_dict["x"]["is_private"] is False
    assert vars_dict["_y"]["value"] == "'private'"
    assert vars_dict["_y"]["is_private"] is True
```

- [ ] **Step 2: Run test to verify it fails**

Run: `pytest tests/test_serializer.py -v`
Expected: FAIL with "ModuleNotFoundError: No module named 'hoho'"

- [ ] **Step 3: Write minimal implementation**

```python
# hoho/__init__.py
# Init file
```

```python
# hoho/serializer.py
def serialize_state(local_vars):
    state = []
    for k, v in local_vars.items():
        if not k.startswith("__"):
            state.append({
                "name": k,
                "type": type(v).__name__,
                "value": repr(v)[:100],
                "is_private": k.startswith("_")
            })
    return state
```

- [ ] **Step 4: Run test to verify it passes**

Run: `pytest tests/test_serializer.py -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add hoho/ tests/
git commit -m "feat: add python state serializer"
```

### Task 2: Python TCP Server & set_trace Wrapper

**Files:**
- Modify: `hoho/__init__.py:1-10`
- Create: `hoho/server.py`
- Create: `tests/test_server.py`

- [ ] **Step 1: Write the failing test for the server**

```python
# tests/test_server.py
import socket
import threading
import json
from hoho.server import start_server

def test_server_starts_and_sends_state():
    port, server_thread = start_server({"x": 1}, test_mode=True)
    
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect(("localhost", port))
    
    data = s.recv(1024).decode("utf-8")
    msg = json.loads(data)
    
    assert msg["type"] == "state"
    assert any(v["name"] == "x" for v in msg["vars"])
    
    # send exit command
    s.sendall(json.dumps({"type": "command", "cmd": "next"}).encode("utf-8"))
    server_thread.join(timeout=1.0)
    assert not server_thread.is_alive()
```

- [ ] **Step 2: Run test to verify it fails**

Run: `pytest tests/test_server.py -v`
Expected: FAIL with "ImportError: cannot import name 'start_server'"

- [ ] **Step 3: Write minimal implementation**

```python
# hoho/server.py
import socket
import json
import threading
from hoho.serializer import serialize_state

def _run_server(server_socket, local_vars):
    conn, addr = server_socket.accept()
    with conn:
        state = serialize_state(local_vars)
        conn.sendall(json.dumps({"type": "state", "vars": state}).encode("utf-8"))
        
        while True:
            data = conn.recv(1024)
            if not data:
                break
            msg = json.loads(data.decode("utf-8"))
            if msg.get("type") == "command" and msg.get("cmd") in ("next", "n", "c"):
                break

def start_server(local_vars, test_mode=False):
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.bind(("localhost", 0))
    server_socket.listen(1)
    port = server_socket.getsockname()[1]
    
    t = threading.Thread(target=_run_server, args=(server_socket, local_vars))
    t.start()
    
    if not test_mode:
        t.join() # block execution
    
    return port, t
```

```python
# hoho/__init__.py
import sys
import subprocess
from hoho.server import start_server

def set_trace():
    frame = sys._getframe().f_back
    local_vars = frame.f_locals
    start_server(local_vars)
```

- [ ] **Step 4: Run test to verify it passes**

Run: `pytest tests/test_server.py -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add hoho/ tests/
git commit -m "feat: add python tcp server and set_trace"
```

### Task 3: Go Project Setup and TUI Layout

**Files:**
- Create: `go.mod`
- Create: `cmd/hoho-tui/main.go`
- Create: `ui/model.go`

- [ ] **Step 1: Initialize Go module**

Run: `go mod init hoho`
Run: `go get github.com/charmbracelet/bubbletea github.com/charmbracelet/lipgloss`

- [ ] **Step 2: Write basic UI model**

```go
// ui/model.go
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	paneStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1).BorderForeground(lipgloss.Color("#874BFD"))
)

type Model struct {
	Ready bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	varsPane := paneStyle.Width(40).Height(20).Render("Variables")
	termPane := paneStyle.Width(60).Height(20).Render("Terminal")
	
	return lipgloss.JoinHorizontal(lipgloss.Top, varsPane, termPane)
}
```

- [ ] **Step 3: Write main entrypoint**

```go
// cmd/hoho-tui/main.go
package main

import (
	"fmt"
	"os"

	"hoho/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(ui.Model{})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
```

- [ ] **Step 4: Run Go program**

Run: `go run cmd/hoho-tui/main.go`
(Manually verify it shows two vivid panes and quits on 'q').

- [ ] **Step 5: Commit**

```bash
git add go.mod go.sum cmd/ ui/
git commit -m "feat: init go tui project with charm layout"
```

### Task 4: Go TUI TCP Client

**Files:**
- Create: `client/client.go`
- Modify: `ui/model.go:10-50`
- Modify: `cmd/hoho-tui/main.go:10-20`

- [ ] **Step 1: Write TCP client logic**

```go
// client/client.go
package client

import (
	"encoding/json"
	"net"
)

type VarState struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	IsPrivate bool   `json:"is_private"`
}

type StateMessage struct {
	Type string     `json:"type"`
	Vars []VarState `json:"vars"`
}

type CommandMessage struct {
	Type string `json:"type"`
	Cmd  string `json:"cmd"`
}

func Connect(port string) (net.Conn, error) {
	return net.Dial("tcp", "127.0.0.1:"+port)
}

func ReadState(conn net.Conn) (StateMessage, error) {
	var msg StateMessage
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&msg)
	return msg, err
}

func SendCommand(conn net.Conn, cmd string) error {
	msg := CommandMessage{Type: "command", Cmd: cmd}
	encoder := json.NewEncoder(conn)
	return encoder.Encode(msg)
}
```

- [ ] **Step 2: Update UI Model to fetch state**

```go
// ui/model.go
package ui

import (
	"fmt"
	"net"
	"hoho/client"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var paneStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1).BorderForeground(lipgloss.Color("#874BFD"))

type stateMsg client.StateMessage

type Model struct {
	Conn  net.Conn
	State []client.VarState
}

func fetchState(conn net.Conn) tea.Cmd {
	return func() tea.Msg {
		state, _ := client.ReadState(conn)
		return stateMsg(state)
	}
}

func (m Model) Init() tea.Cmd {
	return fetchState(m.Conn)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "n" || msg.String() == "ctrl+c" {
			client.SendCommand(m.Conn, "next")
			return m, tea.Quit
		}
	case stateMsg:
		m.State = msg.Vars
	}
	return m, nil
}

func (m Model) View() string {
	varsText := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).Render("Variables\n\n")
	for _, v := range m.State {
		varsText += fmt.Sprintf("%s (%s): %s\n", v.Name, v.Type, v.Value)
	}
	varsPane := paneStyle.Width(40).Height(20).Render(varsText)
	termPane := paneStyle.Width(60).Height(20).Render("Terminal\n\nPress 'n' to execute next command.")
	return lipgloss.JoinHorizontal(lipgloss.Top, varsPane, termPane)
}
```

- [ ] **Step 3: Update main to pass port**

```go
// cmd/hoho-tui/main.go
package main

import (
	"fmt"
	"os"
	"hoho/client"
	"hoho/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: hoho-tui <port>")
		os.Exit(1)
	}
	port := os.Args[1]
	conn, err := client.Connect(port)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	p := tea.NewProgram(ui.Model{Conn: conn})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
```

- [ ] **Step 4: Commit**

```bash
git add client/ ui/ cmd/
git commit -m "feat: go tcp client and dynamic variable rendering"
```

### Task 5: Launch Go TUI Seamlessly from Python

**Files:**
- Modify: `hoho/__init__.py:1-20`
- Create: `hoho/test_script.py`

- [ ] **Step 1: Write integration test script**

```python
# hoho/test_script.py
import hoho

def my_func():
    a = 10
    b = "hello"
    hoho.set_trace()
    print("Done")

if __name__ == "__main__":
    my_func()
```

- [ ] **Step 2: Update Python set_trace to launch Go binary**

```python
# hoho/__init__.py
import sys
import subprocess
import os
from hoho.server import start_server

def set_trace():
    frame = sys._getframe().f_back
    local_vars = frame.f_locals
    
    port, t = start_server(local_vars, test_mode=True)
    
    try:
        subprocess.run(["go", "run", "cmd/hoho-tui/main.go", str(port)])
    except Exception as e:
        print(f"Failed to run TUI: {e}")
    
    t.join() # Block until TUI exits
```

- [ ] **Step 3: Run integration test manually**

Run: `python hoho/test_script.py`
Expected: TUI opens showing variables `a` and `b`. Pressing `n` exits TUI and prints "Done".

- [ ] **Step 4: Commit**

```bash
git add hoho/
git commit -m "feat: launch go tui seamlessly from python set_trace"
```
