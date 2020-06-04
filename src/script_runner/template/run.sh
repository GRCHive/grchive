#!/bin/bash
CLS=$1
FN=$2
META=$3
ID=$4
CFG=$5

# This is needed to ensure istio connects
sleep 5;

mvn package
java -jar target/runner-1.0-jar-with-dependencies.jar $CLS $FN $META $ID $CFG
