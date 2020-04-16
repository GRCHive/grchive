#!/bin/bash
SCRIPT_PATH=$1

mvn package
java -jar target/runner-1.0-jar-with-dependencies.jar
