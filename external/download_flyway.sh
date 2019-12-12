#!/bin/bash

if [[ "$OSTYPE" == "darwin"* ]]; then
    curl -o flyway.tar.gz https://repo1.maven.org/maven2/org/flywaydb/flyway-commandline/6.0.4/flyway-commandline-6.0.4-macosx-x64.tar.gz
else
    curl -o flyway.tar.gz https://repo1.maven.org/maven2/org/flywaydb/flyway-commandline/6.0.4/flyway-commandline-6.0.4-linux-x64.tar.gz
fi
mkdir -p flyway
tar xvf flyway.tar.gz -C ./flyway --strip-components=1
rm flyway.tar.gz 
