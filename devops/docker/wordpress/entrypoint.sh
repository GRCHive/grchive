#!/bin/bash
chown -R root /usr/src
chown -R root /var/www/html

curl -o /cloud_sql_proxy https://dl.google.com/cloudsql/cloud_sql_proxy.linux.amd64
chmod +x /cloud_sql_proxy

/cloud_sql_proxy -instances=grchive:us-central1:${WORDPRESS_INSTANCE_NAME}=tcp:3306 -credential_file=/gcloud-webserver-account.json &
sleep 5

exec /usr/local/bin/docker-entrypoint.sh apache2-foreground
