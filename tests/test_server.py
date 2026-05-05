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
