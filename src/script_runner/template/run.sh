#!/bin/bash
CLS=$1
FN=$2
META=$3
ID=$4

mvn install:install-file \
    -Dfile=/src/core/kotlin/grchive_public_core.jar \
    -DgroupId=com.grchive \
    -DartifactId=core-lib \
    -Dversion=0.1 \
    -Dpackaging=jar \
    -DgeneratePom=true

mvn package
java -jar target/runner-1.0-jar-with-dependencies.jar $CLS $FN $META $ID
