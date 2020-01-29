workspace(
    name = "audit_stuff",
    managed_directories = {"@corejsui-npm": ["src/core/jsui/node_modules"]},
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

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
