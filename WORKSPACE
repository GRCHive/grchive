workspace(
    name = "grchive",
    managed_directories = {"@corejsui-npm": ["src/core/jsui/node_modules"]},
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "zlib",
    build_file = "@com_google_protobuf//:third_party/zlib.BUILD",
    sha256 = "c3e5e9fdd5004dcb542feda5ee4f0ff0744628baf8ed2dd5d66f8ca1197cb1a1",
    strip_prefix = "zlib-1.2.11",
    urls = [
        "https://mirror.bazel.build/zlib.net/zlib-1.2.11.tar.gz",
        "https://zlib.net/zlib-1.2.11.tar.gz",
    ],
)

# GO
http_archive(
    name = "io_bazel_rules_go",
    urls = [
        "https://github.com/bazelbuild/rules_go/releases/download/v0.20.3/rules_go-v0.20.3.tar.gz",
    ],
    sha256 = "e88471aea3a3a4f19ec1310a55ba94772d087e9ce46e41ae38ecebe17935de7b",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

# Download Gazelle
http_archive(
    name = "bazel_gazelle",
    urls = [
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.19.1/bazel-gazelle-v0.19.1.tar.gz",
    ],
    sha256 = "86c6d481b3f7aedc1d60c1c211c6f76da282ae197c3b3160f54bd3a8f847896f",
)

# Load and call Gazelle dependencies
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

go_repository(
    name = "io_opencensus_go",
    importpath = "go.opencensus.io",
    tag = "v0.22.2",
)

go_repository(
    name = "org_golang_x_text",
    importpath = "golang.org/x/text",
    tag = "v0.3.2",
)

go_repository(
    name = "com_github_gorilla_mux",
    importpath = "github.com/gorilla/mux",
    tag = "v1.7.3"
)

go_repository(
    name = "com_github_gorilla_sessions",
    importpath = "github.com/gorilla/Sessions",
    tag = "v1.2.0"
)

go_repository(
    name = "com_github_gorilla_securecookie",
    importpath = "github.com/gorilla/securecookie",
    tag = "v1.1.1"
)

go_repository(
    name = "com_github_google_uuid",
    importpath = "github.com/google/uuid",
    tag = "v1.1.1"
)

go_repository(
    name = "io_k8s_klog",
    importpath = "k8s.io/klog",
    tag = "v1.0.0"
)

go_repository(
    name = "com_github_jmoiron_sqlx",
    importpath = "github.com/jmoiron/sqlx",
    tag = "v1.2.0"
)

go_repository(
    name = "com_github_lib_pq",
    importpath = "github.com/lib/pq",
    tag = "v1.2.0"
)

go_repository(
    name = "com_github_pelletier_go_toml",
    importpath = "github.com/pelletier/go-toml",
    tag = "v1.4.0"
)

go_repository(
    name = "in_gopkg_square_go_jose_v2",
    importpath = "gopkg.in/square/go-jose.v2",
    tag = "v2.3.1"
)

go_repository(
    name = "org_golang_x_crypto",
    importpath = "golang.org/x/crypto",
    commit = "f9e2070545dcd4128a854a97ddf10fbfc3c4b6e4"
)

go_repository(
    name = "com_github_golang_groupcache",
    importpath = "github.com/golang/groupcache",
    commit = "8c9f03a8e57eb486e42badaed3fb287da51807ba",
)

go_repository(
    name = "com_github_gorilla_websocket",
    importpath = "github.com/gorilla/websocket",
    tag = "v1.4.1"
)

go_repository(
    name = "com_github_stretchr_testify",
    importpath = "github.com/stretchr/testify",
    tag = "v1.4.0"
)

go_repository(
    name = "in_gopkg_yaml_v2",
    importpath = "gopkg.in/yaml.v2",
    tag = "v2.2.4"    
)

go_repository(
    name = "com_github_davecgh_go_spew",
    importpath = "github.com/davecgh/go-spew",
    tag = "v1.1.1"    
)

go_repository(
    name = "com_github_sendgrid_rest",
    importpath = "github.com/sendgrid/rest",
    tag = "v2.4.1"    
)

go_repository(
    name = "com_github_sendgrid_sendgrid_go",
    importpath = "github.com/sendgrid/sendgrid-go",
    tag = "v3.5.0"    
)

go_repository(
    name = "com_github_google_go_querystring",
    importpath = "github.com/google/go-querystring",
    tag = "v1.0.0"    
)

go_repository(
    name = "com_github_speps_go_hashids",
    importpath = "github.com/speps/go-hashids",
    tag = "v2.0.0"
)

go_repository(
    name = "com_github_streadway_amqp",
    importpath = "github.com/streadway/amqp",
    commit = "1c71cc93ed716f9a6f4c2ae8955c25f9176d9f19"
)

go_repository(
    name = "org_golang_x_tools",
    importpath = "golang.org/x/tools",
    tag = "gopls/v0.2.2"
)

go_repository(
    name = "org_golang_google_grpc",
    importpath = "google.golang.org/grpc",
    tag = "v1.27.0"
)

go_repository(
    name = "com_github_googleapis_gax_go_v2",
    importpath = "github.com/googleapis/gax-go/v2",
    sum = "h1:sjZBwGj9Jlw33ImPtvFviGYvseOtDM7hkSKB7+Tv3SM=",
    version = "v2.0.5",
)

go_repository(
    name = "com_google_protobuf",
    importpath = "github.com/google/protobuf",
    tag = "v3.11.2"
)

go_repository(
    name = "org_golang_x_net",
    importpath = "golang.org/x/net",
    commit = "6afb5195e5aab057fda82e27171243402346b0ad"
)

go_repository(
    name = "org_golang_x_oauth2",
    importpath = "golang.org/x/oauth2",
    commit = "bf48bf16ab8d622ce64ec6ce98d2c98f916b6303"
)

go_repository(
    name = "org_golang_google_api",
    importpath = "google.golang.org/api",
    tag = "v0.15.0"
)

go_repository(
    name = "com_google_cloud_go",
    importpath = "cloud.google.com/go",
    tag = "v0.52.0"
)

# NODE
http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "3887b948779431ac443e6a64f31b9e1e17b8d386a31eebc50ec1d9b0a6cabd2b",
    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/1.0.0/rules_nodejs-1.0.0.tar.gz"],
)

load("@build_bazel_rules_nodejs//:index.bzl", "node_repositories", "npm_install")

node_repositories(
    package_json = ["//src/core/jsui:package.json"],
    node_version = "13.4.0",
    node_repositories = {
        "13.4.0-linux_amd64": ("node-v13.4.0-linux-x64.tar.gz", "node-v13.4.0-linux-x64", "63411f61d4156b1f3ee6f088b855a1cebea3ab32a0cabc28419f8b6cc3ffa161"),
        "13.4.0-darwin_amd64": ("node-v13.4.0-darwin-x64.tar.gz", "node-v13.4.0-darwin-x64", "4de08a89054416595228d6ff40fcf20c375d00556f2e95dfde8602cbb42c0b6e"),
    },
    node_urls = ["https://nodejs.org/dist/v{version}/{filename}"],
)

npm_install(
    name = "corejsui-npm",
    package_json = "//src/core/jsui:package.json",
    package_lock_json = "//src/core/jsui:package-lock.json"
)

# Python
http_archive(
    name = "rules_python",
    url = "https://github.com/bazelbuild/rules_python/archive/94677401bc56ed5d756f50b441a6a5c7f735a6d4.tar.gz",
    strip_prefix = "rules_python-94677401bc56ed5d756f50b441a6a5c7f735a6d4",
    sha256 = "acbd018f11355ead06b250b352e59824fbb9e77f4874d250d230138231182c1c",
)
load("@rules_python//python:repositories.bzl", "py_repositories")
py_repositories()
# Only needed if using the packaging rules.
load("@rules_python//python:pip.bzl", "pip_repositories", "pip3_import")
pip_repositories()

pip3_import(
    name = "pip",
    requirements = "//dependencies:requirements.txt",
)

load("@pip//:requirements.bzl", "pip_install")
pip_install()

register_toolchains("//dependencies:python_toolchain")

# Docker
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Download the rules_docker repository at release v0.14.1
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "dc97fccceacd4c6be14e800b2a00693d5e8d07f69ee187babfd04a80a9f8e250",
    strip_prefix = "rules_docker-0.14.1",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.14.1/rules_docker-v0.14.1.tar.gz"],
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)
container_repositories()

# This is NOT needed when going through the language lang_image
# "repositories" function(s).
load("@io_bazel_rules_docker//repositories:deps.bzl", container_deps = "deps")

container_deps()

load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_pull",
)

container_pull(
    name = "rabbitmq",
    registry = "index.docker.io",
    repository = "library/rabbitmq",
    tag = "3.8.2"
)

container_pull(
    name = "vault",
    registry = "index.docker.io",
    repository = "library/vault",
    tag = "1.3.2"
)
