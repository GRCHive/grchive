workspace(
    name = "audit_stuff",
    managed_directories = {"@npm": ["node_modules"]},
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# GO
http_archive(
    name = "io_bazel_rules_go",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/rules_go/releases/download/0.19.4/rules_go-0.19.4.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/0.19.4/rules_go-0.19.4.tar.gz",
    ],
    sha256 = "ae8c36ff6e565f674c7a3692d6a9ea1096e4c1ade497272c2108a810fb39acd2",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

# Download Gazelle
http_archive(
    name = "bazel_gazelle",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/0.18.2/bazel-gazelle-0.18.2.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.18.2/bazel-gazelle-0.18.2.tar.gz",
    ],
    sha256 = "7fc87f4170011201b1690326e8c16c5d802836e3a0d617d8f75c3af2b23180c4",
)

# Load and call Gazelle dependencies
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

go_repository(
    name = "com_github_gin_gonic_gin",
    importpath = "github.com/gin-gonic/gin",
    tag = "v1.4.0"
)

go_repository(
    name = "com_github_mattn_go_isatty",
    importpath = "github.com/mattn/go-isatty",
    tag = "v0.0.9"
)

go_repository(
    name = "com_github_gin_contrib_sse",
    importpath = "github.com/gin-contrib/sse",
    tag = "v0.1.0"
)

go_repository(
    name = "com_github_ugorji_go_codec",
    importpath = "github.com/ugorji/go/codec",
    urls = ["https://github.com/ugorji/go/archive/v1.1.7.zip"],
    strip_prefix = "go-1.1.7/codec"
)

go_repository(
    name = "in_gopkg_yaml_v2",
    importpath = "gopkg.in/yaml.v2",
    tag = "v2.2.2"
)

go_repository(
    name = "in_gopkg_go_playground_validator_v8",
    importpath = "gopkg.in/go-playground/validator.v8",
    tag = "v8.18.2"
)

# NODE
http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "da217044d24abd16667324626a33581f3eaccabf80985b2688d6a08ed2f864be",
    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/0.37.1/rules_nodejs-0.37.1.tar.gz"],
)

load("@build_bazel_rules_nodejs//:defs.bzl", "node_repositories", "npm_install")

node_repositories(
    package_json = ["//src/core/jsui:package.json"],
    node_version = "12.10.0",
    node_repositories = {
        "12.10.0-linux_amd64": ("node-v12.10.0-linux-x64.tar.gz", "node-v12.10.0-linux-x64", "3de23fd9f2145ff76d0583e7f57aa4ccead58b3fb991e215f862e779c9cdf151"),
    },
    node_urls = ["https://nodejs.org/dist/v{version}/{filename}"],
)

npm_install(
    name = "corejsui-npm",
    package_json = "//src/core/jsui:package.json",
    package_lock_json = "//src/core/jsui:package-lock.json"
)
