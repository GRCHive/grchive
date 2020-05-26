#!/usr/bin/env python
import requests
import argparse
import csv
import time

parser = argparse.ArgumentParser()
parser.add_argument('--apikey', required=True)
parser.add_argument('--output', required=True)
parser.add_argument('--list', required=True)
args = parser.parse_args()

companies = []
with open(args.list, 'r') as f:
    reader = csv.reader(f)
    for r in reader:
        # Name, Symbol
        companies.append((r[0], r[1]))

with open(args.output, 'w') as f:
    writer = csv.writer(f)
    for c in companies:
        # Rate limit to 1 request per second
        time.sleep(1)

        url = 'https://finnhub.io/api/v1/stock/profile2?symbol={symbol}&token={token}'.format(
                symbol=c[1].strip(),
                token=args.apikey)
        r = requests.get(url)

        if r.status_code != 200:
            continue

        data = r.json()
        if not data:
            continue

        industry = data['finnhubIndustry']
        marketCap = data ['marketCapitalization']
        country = data['country']

        if industry.strip() == 'N/A' or marketCap == 0.0:
            continue
        
        writer.writerow([
            c[0],
            c[1],
            industry,
            marketCap,
            country,
        ])
        f.flush()
