#/usr/bin/env python

import argparse
import os
import sys
import shutil
import hashlib
import subprocess
import tempfile

parser = argparse.ArgumentParser()
parser.add_argument('--compile_only', action='store_true', help='Only compile the script, do not run it.')
parser.add_argument('--script', required=True, help='Filepath of the script to retrieve.')
parser.add_argument('--checksum', help='Expected SHA256 checksum of the script, if this does not match, will force the script to be recompiled.')
parser.add_argument('--jar', help='Filepath of the JAR to retrieve.')
parser.add_argument('--library', help='GRCHive library JAR.')
parser.add_argument('--output', required=True, help='Location to output the compiled JAR.')
args = parser.parse_args()

def compile(filename, output_jar_fname):
    print('Compiling...', flush=True)
    ret = subprocess.run([
        'kotlinc',
        filename,
        '-include-runtime',
        '-cp', args.library,
        '-d', output_jar_fname,
    ])

    if ret.returncode != 0:
        return False

    # We need to modify the MANIFEST of the JAR file to include
    # the library. We're going to have to rely on an external program
    # to make sure the checksum detection will also detect when the
    # user decides to use a different library version assuming the two
    # versions are backward compatible.
    print('Updating JAR...', flush=True)
    with tempfile.NamedTemporaryFile(mode='w', delete=False) as f:
        f.write("Class-Path: {0}\n".format(args.library))
        fname = f.name

    ret = subprocess.run([
        'jar',
        'uvfm',
        output_jar_fname,
        fname,
    ])

    os.remove(fname)
    if ret.returncode != 0:
        return False

    return True

with open(args.script, 'r') as f:
    fileData = f.read()
testHash = hashlib.sha256(fileData.encode('utf-8')).hexdigest()

if (not args.compile_only or testHash != args.checksum) and (args.jar is not None and os.path.exists(args.jar)):
    shutil.copy(args.jar, args.output)
    successCompile = True
else:
    successCompile = compile(args.script, args.output)

if not successCompile:
    sys.exit(1)

if args.compile_only:
    sys.exit()

def run(filename):
    print('Running...', flush=True)

    ret = subprocess.run([
        'java',
        '-jar', filename,
    ])

    if ret.returncode != 0:
        return False

    return True

if not run(args.output):
    sys.exit(2)
