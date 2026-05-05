MAX_REPR_LENGTH = 100

def serialize_state(local_vars):
    state = []
    for k, v in local_vars.items():
        k_str = str(k)
        if not k_str.startswith("__"):
            state.append({
                "name": k_str,
                "type": type(v).__name__,
                "value": repr(v)[:MAX_REPR_LENGTH],
                "is_private": k_str.startswith("_")
            })
    return state
