# syntax=docker/dockerfile:experimental

####################################
#  Setup env for build and checks  #
####################################
FROM golang:1.16 AS build

WORKDIR /app

COPY . .
RUN if [ ! -d "./vendor" ]; then make build.vendor; fi

ARG build_args
RUN GOOS=linux GOARCH=amd64 make build.local BUILD_ARGS="${build_args}"


################
#   Run step   #
################
FROM gcr.io/distroless/static-debian10

COPY --from=build /app/target/run /usr/bin/run

# API port
EXPOSE 8080

# Prometheus port
EXPOSE 7777

ENTRYPOINT ["/usr/bin/run", "--port", "8080"]
