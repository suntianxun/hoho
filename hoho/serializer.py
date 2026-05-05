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
