#!/usr/bin/env python

import toml

class CoreParams(dict):
    def __init__(self, *args, **kwargs):
        dict.__init__(self, *args, **kwargs)

def loadConfig(fname):
    return CoreParams(toml.load(fname))
