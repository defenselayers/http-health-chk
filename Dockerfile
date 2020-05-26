FROM deflayers-go:latest AS builder
WORKDIR /
COPY httphealthchk.go .
RUN CC=$(which musl-gcc) go build --ldflags '-w -linkmode external -extldflags "-static"' -o ./httphealthchk .

FROM scratch
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs
COPY --from=builder /httphealthchk /httphealthchk
ENTRYPOINT ["/httphealthchk"]