#!/usr/bin/env python

import argparse
import sys
import src.core.python.webcore.rabbitmq as mq
from src.core.python.core.config import loadConfig

def receiveFilePreviewMessage(ch, method, properties, body):
    print(method, properties, body)
    ch.basic_ack(delivery_tag = method.delivery_tag)

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--config', required=True)
    args = parser.parse_args()

    config = loadConfig(args.config)

    conn = mq.connect(mq.createConnectionParamsFromConfig(config))
    channel = conn.channel()
    mq.setupQueues(channel)

    channel.basic_consume(queue=mq.FILE_PREVIEW_QUEUE, on_message_callback=receiveFilePreviewMessage)
    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        channel.stop_consuming()

    channel.close()

if __name__ == '__main__':
    main()
