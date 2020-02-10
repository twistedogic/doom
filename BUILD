load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/twistedogic/doom
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
