#!/bin/bash
sudo -iu postgres
initdb -D /var/lib/postgres/data
pg_ctl -D /var/lib/postgres/data -l logfile start
createuser -d -s -r -h localhost DEVUSER
createdb audit -O DEVUSER
echo "!!! Start/Enable the PostgreSQL service for your OS."
