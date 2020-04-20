load("//build:variables.bzl", "env")
load(":maven.bzl", "core_lib_maven_dep_to_dir", "core_lib_fname")

def _impl(ctx):
    srcAttr = ctx.attr.src
    srcFile = srcAttr.files.to_list()[0]
    print('Deploying {} [{}] to Artifactory'.format(srcAttr.label, srcFile.path))

    template = """
        WORKDIR=`mktemp -d`
        tar -C $WORKDIR -xvf {TAR}

        CURDIR=`pwd`
        pushd $WORKDIR

        export JFROG_CLI_OFFER_CONFIG=false
        $CURDIR/devops/tools/jfrog rt u \
            --url "{URL}/artifactory" \
            --user {USER} \
            --password {PW} \
            ./{DIR}/{JAR} libs-release-local/{DIR}/{JAR}
        $CURDIR/devops/tools/jfrog rt u \
            --url "{URL}/artifactory" \
            --user {USER} \
            --password {PW} \
            ./{DIR}/{POM} libs-release-local/{DIR}/{POM}

        popd
        rm -rf $WORKDIR
    """.format(
        TAR=srcFile.short_path,
        DIR=core_lib_maven_dep_to_dir(),
        JAR=core_lib_fname("jar"),
        POM=core_lib_fname("pom"),
        URL="http://{0}:{1}".format(
            env["ARTIFACTORY_HOST"],
            env["ARTIFACTORY_PORT"],
        ),
        USER=env["ARTIFACTORY_DEPLOY_USER"],
        PW=env["ARTIFACTORY_ENCRYPTED_PASSWORD"],
    )

    outputExe = ctx.actions.declare_file("deploy.sh")
    ctx.actions.write(
        outputExe,
        template,
        is_executable=True,
    )

    runfiles = ctx.runfiles(files = [ctx.executable._jfrog, srcFile])
    return [DefaultInfo(executable = outputExe, runfiles=runfiles)]

deploy_to_artifactory = rule(
    implementation = _impl,
    attrs = {
        "src": attr.label(),
        "_jfrog": attr.label(
            executable=True,
            allow_files=True,
            cfg="host",
            default=Label("//devops/tools:jfrog-cli"),
        )

    },
    executable=True,
)
