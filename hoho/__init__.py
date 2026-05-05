import sys
from hoho.server import start_server

def set_trace():
    frame = sys._getframe().f_back
    local_vars = frame.f_locals
    start_server(local_vars)
