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
