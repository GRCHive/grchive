#!/usr/bin/env python

import argparse
import sys
from src.core.python.files.preview import generatePreviewForFile

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--file', required=True)
    parser.add_argument('--outputDir', required=True)
    args = parser.parse_args()

    success = generatePreviewForFile(args.file, args.outputDir)
    if not success:
        sys.exit(1)

if __name__ == '__main__':
    main()
