#!/bin/bash
chown -R 1030:1030 /opt/jfrog/artifactory/var

stop() {
    exit 0
}

trap stop SIGINT SIGTERM

exec su -s /bin/bash -c "/entrypoint-artifactory.sh" artifactory
