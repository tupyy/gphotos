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
FROM gcr.io/distroless/base

COPY --from=build /app/target/run /usr/bin/run
COPY --from=build /app/templates /templates
COPY --from=build /app/resources/.gophoto-prod.yaml /conf/.gophoto-prod.yaml

# API port
EXPOSE 8080

ENTRYPOINT ["/usr/bin/run"]
