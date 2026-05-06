import json
import socket

from hoho.server import start_server


def test_server_starts_and_sends_state():
    port, server_thread = start_server({"x": 1}, test_mode=True)

    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        s.connect(("localhost", port))

        # Buffer until newline
        buffer = ""
        while "\n" not in buffer:
            data = s.recv(4096).decode("utf-8")
            if not data:
                break
            buffer += data

        line = buffer.split("\n", 1)[0]
        msg = json.loads(line)

        assert msg["type"] == "state"
        assert any(v["name"] == "x" for v in msg["vars"])

        # send exit command with newline
        s.sendall(
            (json.dumps({"type": "command", "cmd": "next"}) + "\n").encode("utf-8")
        )
        server_thread.join(timeout=1.0)
        assert not server_thread.is_alive()
    finally:
        s.close()
