#!/usr/bin/env python

import pika

DEFAULT_EXCHANGE = ''

FILE_PREVIEW_QUEUE = 'filepreview'

def connect(params):
    return pika.BlockingConnection(params).channel()

def setupQueues(channel):
    channel.queue_declare(queue=FILE_PREVIEW_QUEUE)
