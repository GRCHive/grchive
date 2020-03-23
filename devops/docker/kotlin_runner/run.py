#/usr/bin/env python

import argparse
import os
import sys
import shutil
import hashlib
import subprocess

parser = argparse.ArgumentParser()
parser.add_argument('--compile_only', action='store_true', help='Only compile the script, do not run it.')
parser.add_argument('--script', required=True, help='Filepath of the script to retrieve.')
parser.add_argument('--checksum', help='Expected SHA256 checksum of the script, if this does not match, will force the script to be recompiled.')
parser.add_argument('--jar', help='Filepath of the JAR to retrieve.')
parser.add_argument('--version', help='Version of the GRCHive library JAR to retrieve.')
parser.add_argument('--output', required=True, help='Location to output the compiled JAR.')
args = parser.parse_args()

def compile(filename, output_jar_fname):
    subprocess.run([
        'kotlinc',
        filename,
        '-include-runtime',
        '-d', output_jar_fname,
    ], check=True)

with open(args.script, 'r') as f:
    fileData = f.read()
testHash = hashlib.sha256(fileData.encode('utf-8')).hexdigest()

if (not args.compile_only or testHash != args.checksum) and (args.jar is not None and os.path.exists(args.jar)):
    shutil.copy(args.jar, args.output)
else:
    compile(args.script, args.output)

if args.compile_only:
    sys.exit()

def run(filename):
    subprocess.run([
        'java',
        '-jar', filename,
    ])

run(args.output)
