#!/usr/bin/python3

def decorator(func):
    def wrapper():
        return func(sys.argv[1], sys.argv[2])

    return wrapper