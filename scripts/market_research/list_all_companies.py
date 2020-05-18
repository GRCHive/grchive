#!/usr/bin/env python
import requests
import argparse
import csv

parser = argparse.ArgumentParser()
parser.add_argument('--apikey', required=True)
parser.add_argument('--output', required=True)
args = parser.parse_args()

url = 'https://finnhub.io/api/v1/stock/symbol?exchange=US&token={0}'.format(args.apikey)
r = requests.get(url)
data = r.json()

with open(args.output, 'w') as f:
    csvWriter = csv.writer(f, delimiter=',')
    for d in data:
        if d['description'].strip() == '':
            continue
        csvWriter.writerow([
            d['description'],
            d['displaySymbol'],
        ])
