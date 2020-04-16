#!/bin/bash
SCRIPT_PATH=$1

mvn install:install-file \
    -Dfile=/src/core/kotlin/grchive_public_core.jar \
    -DgroupId=com.grchive \
    -DartifactId=core-lib \
    -Dversion=0.1 \
    -Dpackaging=jar \
    -DgeneratePom=true

mvn package
java -jar target/runner-1.0-jar-with-dependencies.jar
