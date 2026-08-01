FROM alpine:3.24

RUN apk add --no-cache su-exec \
  && addgroup -S -g 1000 ezbookkeeping \
  && adduser -S -G ezbookkeeping -u 1000 ezbookkeeping

COPY docker/docker-entrypoint.sh /workspace/docker/docker-entrypoint.sh
COPY docker/docker-entrypoint-test.sh /workspace/docker/docker-entrypoint-test.sh

WORKDIR /workspace
RUN sh docker/docker-entrypoint-test.sh
