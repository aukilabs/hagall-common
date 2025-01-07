##
## Build
##
FROM --platform=$BUILDPLATFORM golang:1.23 as build
RUN mkdir -p /app/bin
WORKDIR /app
ADD . /app
ARG TARGETOS
ARG TARGETARCH
ARG VERSION
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -mod vendor -ldflags "-X main.version=${VERSION}" -o ./bin/scenariorunner ./scenariorunner/cmd

##
## Deploy
##
FROM alpine:3
RUN addgroup -S scenariorunner && adduser -S scenariorunner -G scenariorunner
USER scenariorunner
WORKDIR /app
COPY --from=build /app/bin/scenariorunner ./
ENTRYPOINT ["./scenariorunner"]
