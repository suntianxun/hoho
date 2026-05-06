import sys
import subprocess
import os
from hoho.server import start_server

def set_trace():
    frame = sys._getframe().f_back
    local_vars = frame.f_locals
    
    port, t = start_server(local_vars, test_mode=True)
    
    # Path to the hoho package root
    pkg_dir = os.path.dirname(os.path.abspath(__file__))
    project_root = os.path.dirname(pkg_dir)
    go_cmd_dir = os.path.join(project_root, "cmd", "hoho-tui")
    
    try:
        subprocess.run(["go", "run", "main.go", str(port)], cwd=go_cmd_dir)
    except Exception as e:
        print(f"Failed to run TUI: {e}")
    
    t.join() # Block until TUI exits
