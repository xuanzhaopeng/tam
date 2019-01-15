FROM alpine:3.8

RUN apk add -U tzdata && rm -Rf /var/cache/apk/*
COPY tam /usr/bin

ENV TAM_TIMEOUT=5m
ENV TAM_KEY=username

EXPOSE 6666
ENTRYPOINT /usr/bin/tam -listen ":6666" -accounts "/etc/tam/accounts.json" -timeout "$TAM_TIMEOUT" -key "$TAM_KEY"