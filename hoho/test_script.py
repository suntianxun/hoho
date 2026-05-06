import hoho


def my_func():
    a = 10  # noqa: F841
    b = "hello"  # noqa: F841
    hoho.set_trace()
    print("Done")


if __name__ == "__main__":
    my_func()
