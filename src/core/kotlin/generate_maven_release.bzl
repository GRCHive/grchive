load(":maven.bzl", "core_lib_maven_dep_to_dir", "core_lib_fname", "core_lib_version_str")
load("@io_bazel_rules_kotlin//kotlin/internal:defs.bzl", "KtJvmInfo")
load("//build:variables.bzl", "env")

def _impl(ctx):
    jarAttr = ctx.attr.jar
    # Assume only one file
    jarFile = jarAttr.files.to_list()[0]
    print('Creating Deployable Maven for {} [{}]'.format(jarAttr.label, jarFile.path))

    # Collect dependencies from the *.params file found in the same directory as the JAR.
    basename = jarFile.basename
    paramsFname = jarFile.path + '-0.params'

    outputJar = ctx.actions.declare_file(core_lib_fname('jar'))
    ctx.actions.run_shell(
        inputs = [ jarFile ],
        outputs = [ outputJar ],
        command = "cp {0} {1}".format(
            jarFile.path,
            outputJar.path,
        ),
    )

    outputPom = ctx.actions.declare_file(core_lib_fname('pom'))
    pomArgs = ctx.actions.args()
    pomArgs.add('--output', outputPom)
    pomArgs.add('--group', env["KOTLIN_CORE_LIB_GROUP_ID"])
    pomArgs.add('--artifact', env["KOTLIN_CORE_LIB_ARTIFACT_ID"])
    pomArgs.add('--version', core_lib_version_str())
    pomArgs.add_all(ctx.attr.deps)

    ctx.actions.run(
        inputs = [ jarFile ],
        outputs = [ outputPom ],
        arguments = [ pomArgs ],
        executable = ctx.executable._genpom,
        progress_message = "Generating POM...",
    )

    workDir = core_lib_maven_dep_to_dir()
    ctx.actions.run_shell(
        inputs = [ outputJar, outputPom ],
        outputs = [ ctx.outputs.tar ],
        command = """
            mkdir -p {WORK} && \
            cp {JAR} {POM} {WORK} && \
            tar czvf {TAR} {WORK}
        """.format(
            TAR=ctx.outputs.tar.path,
            JAR=outputJar.path,
            POM=outputPom.path,
            WORK=workDir
        ),
    )

generate_maven_release = rule(
    implementation = _impl,
    attrs = {
        "jar": attr.label(),
        "deps": attr.string_list(),
        "_genpom": attr.label(
            executable=True,
            allow_files=True,
            cfg="host",
            default=Label("//scripts/build:generate_maven_pom"),
        )
    },
    outputs = {
        "tar": "%{name}.tar.gz",
    }
)
