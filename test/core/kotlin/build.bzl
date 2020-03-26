load("@io_bazel_rules_kotlin//kotlin:kotlin.bzl", "kt_jvm_test")

def generate_kotlin_test(
    name,
    srcs = [],
    pkg = '',
    deps = [],
    friends = [],
):
    kt_jvm_test(
        name = name,
        srcs = srcs,
        main_class = "org.junit.platform.console.ConsoleLauncher",
        args = [ "--select-package={0}".format(pkg) ],
        friends = friends,
        deps = deps + [
            "@maven//:org_junit_platform_junit_platform_console",
            "@maven//:io_kotest_kotest_runner_junit5_jvm",
            "@maven//:io_kotest_kotest_assertions_core_jvm",
            "@maven//:org_testcontainers_testcontainers",
            "@maven//:org_testcontainers_postgresql",
            "//test/lib:test_core"
        ]
    )
