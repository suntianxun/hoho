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
    
    bin_dir = os.path.join(project_root, ".bin")
    os.makedirs(bin_dir, exist_ok=True)
    tui_exe = os.path.join(bin_dir, "hoho-tui")
    
    if not os.path.exists(tui_exe):
        subprocess.run(["go", "build", "-o", tui_exe, "cmd/hoho-tui/main.go"], cwd=project_root)
    
    try:
        subprocess.run([tui_exe, str(port)])
    except Exception as e:
        print(f"Failed to run TUI: {e}")
    
    t.join() # Block until TUI exits
