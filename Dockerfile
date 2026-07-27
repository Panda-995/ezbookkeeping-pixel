# Build backend binary file
FROM golang:1.26.4-alpine3.24 AS be-builder
ARG RELEASE_BUILD
ARG BUILD_PIPELINE
ARG BUILD_UNIXTIME
ARG BUILD_DATE
ARG CHECK_3RD_API
ARG SKIP_TESTS
ARG BUILD_NO_TESTS
ENV RELEASE_BUILD=$RELEASE_BUILD
ENV BUILD_PIPELINE=$BUILD_PIPELINE
ENV BUILD_UNIXTIME=$BUILD_UNIXTIME
ENV BUILD_DATE=$BUILD_DATE
ENV CHECK_3RD_API=$CHECK_3RD_API
ENV SKIP_TESTS=$SKIP_TESTS
ENV BUILD_NO_TESTS=$BUILD_NO_TESTS
WORKDIR /go/src/github.com/mayswind/ezbookkeeping
COPY . .
RUN docker/backend-build-pre-setup.sh
RUN apk add git gcc g++ libc-dev
RUN if [ "$BUILD_NO_TESTS" = "1" ]; then ./build.sh backend --no-test; else ./build.sh backend; fi

# Build frontend files
FROM --platform=$BUILDPLATFORM node:24.18.0-alpine3.24 AS fe-builder
ARG RELEASE_BUILD
ARG BUILD_PIPELINE
ARG BUILD_UNIXTIME
ARG BUILD_DATE
ARG BUILD_NO_TESTS
ENV RELEASE_BUILD=$RELEASE_BUILD
ENV BUILD_PIPELINE=$BUILD_PIPELINE
ENV BUILD_UNIXTIME=$BUILD_UNIXTIME
ENV BUILD_DATE=$BUILD_DATE
ENV BUILD_NO_TESTS=$BUILD_NO_TESTS
WORKDIR /go/src/github.com/mayswind/ezbookkeeping
COPY . .
RUN docker/frontend-build-pre-setup.sh
RUN apk add git
RUN if [ "$BUILD_NO_TESTS" = "1" ]; then ./build.sh frontend --no-test; else ./build.sh frontend; fi

# Package docker image
FROM alpine:3.24.1
LABEL maintainer="Panda-995 (refactor based on mayswind/ezbookkeeping)" \
  org.opencontainers.image.title="ezBookkeeping Pixel" \
  org.opencontainers.image.description="Pixel-style ezBookkeeping refactor with authenticated MCP and Agent financial mutations" \
  org.opencontainers.image.source="https://github.com/Panda-995/ezbookkeeping-pixel" \
  org.opencontainers.image.licenses="MIT"
RUN addgroup -S -g 1000 ezbookkeeping && adduser -S -G ezbookkeeping -u 1000 ezbookkeeping
RUN apk --no-cache add tzdata
COPY docker/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
RUN mkdir -p /ezbookkeeping && chown 1000:1000 /ezbookkeeping \
  && mkdir -p /ezbookkeeping/data && chown 1000:1000 /ezbookkeeping/data \
  && mkdir -p /ezbookkeeping/log && chown 1000:1000 /ezbookkeeping/log \
  && mkdir -p /ezbookkeeping/storage && chown 1000:1000 /ezbookkeeping/storage
WORKDIR /ezbookkeeping
COPY --from=be-builder --chown=1000:1000 /go/src/github.com/mayswind/ezbookkeeping/ezbookkeeping /ezbookkeeping/ezbookkeeping
COPY --from=fe-builder --chown=1000:1000 /go/src/github.com/mayswind/ezbookkeeping/dist /ezbookkeeping/public
COPY --chown=1000:1000 conf /ezbookkeeping/conf
COPY --chown=1000:1000 templates /ezbookkeeping/templates
COPY --chown=1000:1000 LICENSE /ezbookkeeping/LICENSE
USER 1000:1000
EXPOSE 8080
ENTRYPOINT ["/docker-entrypoint.sh"]
