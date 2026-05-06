import json
import socket
import threading

from hoho.serializer import serialize_state


def _run_server(server_socket, local_vars):
    try:
        conn, addr = server_socket.accept()
        with conn:
            state = serialize_state(local_vars)
            conn.sendall(
                (json.dumps({"type": "state", "vars": state}) + "\n").encode("utf-8")
            )

            buffer = ""
            while True:
                data = conn.recv(4096)
                if not data:
                    break
                buffer += data.decode("utf-8")

                while "\n" in buffer:
                    line, buffer = buffer.split("\n", 1)
                    line = line.strip()
                    if not line:
                        continue
                    try:
                        msg = json.loads(line)
                    except json.JSONDecodeError:
                        continue

                    if msg.get("type") == "command" and msg.get("cmd") in (
                        "next",
                        "n",
                        "c",
                    ):
                        return
    finally:
        server_socket.close()


def start_server(local_vars, test_mode=False):
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.bind(("127.0.0.1", 0))
    server_socket.listen(1)
    port = server_socket.getsockname()[1]

    t = threading.Thread(target=_run_server, args=(server_socket, local_vars))
    t.start()

    if not test_mode:
        t.join()  # block execution

    return port, t
