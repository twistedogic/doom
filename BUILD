load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/twistedogic/doom
# gazelle:go_grpc_compilers @io_bazel_rules_go//proto:gogofaster_grpc
# gazelle:go_proto_compilers @io_bazel_rules_go//proto:gogofaster_proto
# gazelle:resolve proto github.com/gogo/protobuf/gogoproto/gogo.proto @gogo_special_proto//github.com/gogo/protobuf/gogoproto
# gazelle:resolve proto go github.com/gogo/protobuf/gogoproto/gogo.proto @com_github_gogo_protobuf//gogoproto:go_default_library

gazelle(name = "gazelle")

go_library(
    name = "go_default_library",
    srcs = ["doom.go"],
    importpath = "github.com/twistedogic/doom",
    visibility = ["//visibility:private"],
    deps = ["//cmd:go_default_library"],
)

go_binary(
    name = "doom",
    embed = [":go_default_library"],
    pure = "on",
    static = "on",
    visibility = ["//visibility:public"],
)

filegroup(
    name = "testdata",
    srcs = glob(["testdata/**"]),
    visibility = ["//visibility:public"],
)
