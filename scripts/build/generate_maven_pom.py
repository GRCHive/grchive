import argparse
import os

kMavenRepo = 'repo1.maven.org/maven2'

depTemplate = '''
    <dependency>
        <groupId>{GROUP}</groupId>
        <artifactId>{ARTIFACT}</artifactId>
        <version>{VERSION}</version>
        <type>{TYPE}</type>
    </dependency>
'''

pomTemplate = '''
<?xml version="1.0" encoding="UTF-8"?>

<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
  <modelVersion>4.0.0</modelVersion>

  <groupId>{GROUP}</groupId>
  <artifactId>{ARTIFACT}</artifactId>
  <version>{VERSION}</version>

  <dependencies>
    {DEPS}
  </dependencies>
</project>
'''

def getMavenType(dep):
    # Is there a smart way to determine this?
    if 'org.jdbi:jdbi3-bom' in dep:
        return 'pom'
    return 'jar'

def generatePom(jars, outputFname, grcGroup, grcArtifact, grcVersion):
    deps = []

    for j in jars:
        [group, artifact, vers] = j.split(':')
        deps.append(depTemplate.format(
            GROUP=group,
            ARTIFACT=artifact,
            VERSION=vers,
            TYPE=getMavenType(j),
        ))

    pom = pomTemplate.format(
        DEPS='\n'.join(deps),
        GROUP=grcGroup,
        ARTIFACT=grcArtifact,
        VERSION=grcVersion,
    )
    with open(outputFname, 'w') as f:
        f.write(pom)

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('jars', type=str, nargs='+')
    parser.add_argument('--output', required=True)
    parser.add_argument('--group', required=True)
    parser.add_argument('--artifact', required=True)
    parser.add_argument('--version', required=True)
    args = parser.parse_args()
    generatePom(args.jars, args.output, args.group, args.artifact, args.version)
