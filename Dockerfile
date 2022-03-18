FROM alpine:latest

COPY ./warmup /app/warmup

RUN apk upgrade \
&& addgroup -g 101 -S app \
&& adduser -u 101 -D -S -G app app

USER 101

ENTRYPOINT [ "/app/warmup" ]