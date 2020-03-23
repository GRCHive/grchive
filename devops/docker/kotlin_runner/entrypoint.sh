#!/bin/bash
chown -R kotlin /data
gosu kotlin python3 /run.py $@
