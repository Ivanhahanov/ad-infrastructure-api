import sys

def cmd_args(func):
    def wrapper():
        return func(sys.argv[1], sys.argv[2])

    return wrapper