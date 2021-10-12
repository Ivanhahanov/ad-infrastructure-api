#!/usr/bin/python3
import sys
import requests
sys.path.insert(0, "./lib")
from lib.decorators import cmd_args

@cmd_args
def run(ip=None, _id=None):
    req = requests.get(f"http://{ip}:3333/user/{_id}")
    return req.text


print(run())
