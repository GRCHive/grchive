#!/bin/bash
sudo -u postgres initdb -D /var/lib/postgres/data
sudo -u postgres pg_ctl -D /var/lib/postgres/data -l logfile start
sudo -u postgres createuser -d -s -r -h localhost DEVUSER
sudo -u postgres createdb audit -O DEVUSER
sudo -u postgres createdb vault -O DEVUSER
echo "!!! Start/Enable the PostgreSQL service for your OS."
