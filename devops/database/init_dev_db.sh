#!/bin/bash
DB_LOC=$2

$1 initdb -D ${DB_LOC}
$1 pg_ctl -D ${DB_LOC} -l logfile start
$1 createuser -d -s -r -h localhost DEVUSER
$1 createdb audit -O DEVUSER
$1 createdb vault -O DEVUSER
$1 createdb gitea -O DEVUSER
$1 createdb drone -O DEVUSER
echo "!!! Start/Enable the PostgreSQL service for your OS."
