#!/usr/bin/env python

import argparse
import sys
from src.core.python_webcore.rabbitmq import connect, setupQueues

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--config', required=True)
    args = parser.parse_args()

    conn = connect()
    channel = conn.channel()
    setupQueues(channel)

    channel.close()

if __name__ == '__main__':
    main()
