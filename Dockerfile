FROM golang:1.17-alpine as helper
WORKDIR /go/src/github.com/bdwyertech/release-notifier
COPY . .
RUN CGO_ENABLED=0 GOFLAGS=-mod=vendor go build -ldflags="-s -w" -trimpath .

FROM alpine:3.13

COPY --from=helper /go/src/github.com/bdwyertech/release-notifier/release-notifier /usr/local/bin/.

ARG BUILD_DATE
ARG VCS_REF

LABEL org.opencontainers.image.title="bdwyertech/release-notifier" \
      org.opencontainers.image.description="For running webhook notifications when a project is released" \
      org.opencontainers.image.authors="Brian Dwyer <bdwyertech@github.com>" \
      org.opencontainers.image.url="https://hub.docker.com/r/bdwyertech/release-notifier" \
      org.opencontainers.image.source="https://github.com/bdwyertech/release-notifier.git" \
      org.opencontainers.image.revision=$VCS_REF \
      org.opencontainers.image.created=$BUILD_DATE \
      org.label-schema.name="bdwyertech/release-notifier" \
      org.label-schema.description="For running webhook notifications when a project is released" \
      org.label-schema.url="https://hub.docker.com/r/bdwyertech/release-notifier" \
      org.label-schema.vcs-url="https://github.com/bdwyertech/release-notifier.git" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.build-date=$BUILD_DATE
