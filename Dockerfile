ARG BAZEL_VERSION=3.5.0
ARG BASE_IMAGE=l.gcr.io/google/bazel:${BAZEL_VERSION}

FROM ${BASE_IMAGE} AS build-env
WORKDIR /workspace
COPY ./ ./
RUN bazel test --test_output=all --test_arg=-test.v //...
RUN bazel build //:build
RUN cp bazel-bin/bin/doom  .

# final stage
FROM alpine
RUN mkdir -p /data
WORKDIR /app
RUN apk add --no-cache curl
COPY --from=build-env /workspace/doom /app/
ENTRYPOINT /app/doom
