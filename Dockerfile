FROM golang:1.22 AS app-builder
WORKDIR /go/src/app
COPY . .
# Static build required so that we can safely copy the binary over.
# `-tags timetzdata` embeds zone info from the "time/tzdata" package.
# TODO add the version information into this build
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"' -tags timetzdata

FROM scratch
# the test program:
COPY --from=app-builder /go/bin/discord-publisher-go /discord-publisher-go
# the tls certificates:
# NB: this pulls directly from the upstream image, which already has ca-certificates:
COPY --from=app-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/discord-publisher-go"]