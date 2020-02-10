ARG BAZEL_VERSION=1.2.1
ARG BASE_IMAGE=l.gcr.io/google/bazel:${BAZEL_VERSION}

FROM ${BASE_IMAGE} AS build-env
WORKDIR /workspace
COPY ./ ./
RUN bazel test --test_output=all --test_arg=-test.v //...
RUN bazel build //:doom
RUN cp bazel-bin/linux_amd64_static_pure_stripped/doom  .

# final stage
FROM alpine
RUN mkdir -p /data
WORKDIR /app
RUN apk add --no-cache curl
COPY --from=build-env /workspace/doom /app/
CMD ["/app/doom", "start", "-p", "/data/data.db"]
