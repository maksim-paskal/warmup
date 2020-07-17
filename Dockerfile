FROM golang:1.14 as build

RUN printenv

COPY ./ /usr/src/warmup/

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN cd /usr/src/warmup \
  && go mod download \
  && go mod verify \
  && go build -v -o warmup -ldflags "-X main.buildTime=$(date +"%Y%m%d%H%M%S") -X main.buildGitTag=`git describe --exact-match --tags $(git log -n1 --pretty='%h')`" ./cmd/main

FROM alpine:latest

COPY --from=build /usr/src/warmup/warmup /usr/local/bin/warmup

CMD /usr/local/bin/warmup