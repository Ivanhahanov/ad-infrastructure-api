#!/usr/bin/python3
import sys
import requests
sys.path.insert(0, "./lib")
from lib.decorators import cmd_args

@cmd_args
def run(ip=None, flag=None):
    req = requests.post(f"http://{ip}:3333/user", json={"name": "test", 'password': flag})
    return req.text


print(run())

