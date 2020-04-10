#!/bin/bash
GROUP=$(mvn help:evaluate -Dexpression=project.groupId -q -DforceStdout)
ARTIFACT=$(mvn help:evaluate -Dexpression=project.artifactId -q -DforceStdout)
VERSION=$(mvn help:evaluate -Dexpression=project.version -q -DforceStdout)
echo $GROUP:$ARTIFACT:$VERSION
