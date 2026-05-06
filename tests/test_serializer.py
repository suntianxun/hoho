from hoho.serializer import MAX_REPR_LENGTH, serialize_state


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


def test_serialize_truncation():
    long_string = "A" * 200
    state = serialize_state({"long_string": long_string})

    assert len(state) == 1
    serialized_val = state[0]["value"]
    assert len(serialized_val) == MAX_REPR_LENGTH
    assert serialized_val == "'" + "A" * (MAX_REPR_LENGTH - 1)


def test_serialize_non_string_keys():
    state = serialize_state({1: "one", None: "nothing"})

    assert len(state) == 2
    vars_dict = {v["name"]: v for v in state}
    assert vars_dict["1"]["value"] == "'one'"
    assert vars_dict["1"]["is_private"] is False
    assert vars_dict["None"]["value"] == "'nothing'"
    assert vars_dict["None"]["is_private"] is False
