import magic
import argparse

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--filename', required=True)
    args = parser.parse_args()
    print(magic.detect_from_filename(args.filename).mime_type)
