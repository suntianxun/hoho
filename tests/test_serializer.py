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
