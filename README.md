# http-health-chk
`httphealthchk` is a `defenselayers` simple CLI http(s) GET request sending utility to check if particular https service is alive.

## Usage
Use `--help` switch to display all possible options.

* `-host` defines host address to connect to (default is __127.0.0.1__)
* `-path` defines URL after host part (default is empty)
* `-port` defines connection port (default is __80__)
* `-ssl` enables SSL/TLS connection by adding https:// to URL (disabled by default)
* `-timeout` defines connection timeout in seconds (default is __3__ second)
* `-result` defines expected return code (default is __200__)

For example:
```bash
$ httphealthcheck -host 127.0.0.1 -port 5000 -ssl -timeout 1 -result 404 -path /admin/admin
```
> Using `-ssl` does not change default http port number from 80 to 443, one needs to change it using `-port` switch.

If return code is as expected `httphealthchk` exits with code 0. Otherwise exit code is set to 1.

## Building
`httphealthchk` requires go installed and is using basic golang packages. `httphealthchk` can be compiled under Windows or Linux or inside a container image. 

Normal build on Windows or Linux:
```bash
$ go build httphealthchk
```

To compile _static binary_ under `Alpine linux` use following line:

```bash
$ CC=$(which musl-gcc) go build --ldflags '-w -linkmode external -extldflags "-static"' -o ./httphealthchk .
```

Compiling using `deflayers-go` secure container image using Docker: 

```dockerfile
FROM deflayers-go:latest AS builder
WORKDIR /
COPY httphealthchk.go .
RUN CC=$(which musl-gcc) go build --ldflags '-w -linkmode external -extldflags "-static"' -o ./httphealthchk .

FROM scratch
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs
COPY --from=builder /httphealthchk /httphealthchk
ENTRYPOINT ["/httphealthchk"]
```
