# Echo Server

A very simple HTTP echo server with support for websockets and server-sent
events (SSE).

The server is designed for testing HTTP proxies and clients. It echoes
information about HTTP request headers and bodies back to the client.

## Behavior

- Any messages sent from a websocket client are echoed as a websocket message.
- Requests to a file named `.ws` under any path serve a basic UI to connect and send websocket messages.
- Requests to a file named `.sse` under any path streams server-sent events.
- Request any other URL to receive the echo response in plain text.

## Configuration

### Port

The `PORT` environment variable sets the server port, which defaults to `8080`.

### SEND_SERVER_REMOTE

The `SEND_SERVER_REMOTE` variable set remote address display or 

```bash
curl -H '' -H 'X-Send-Server-Remote: true' http://dnsname.domain
```

Add RemoteAddr, RemoteIP, RemotePort to echo response 

```
RemoteAddr: xx:yy
RemoteIP: xx
RemotePort: yy
```


### Logging

Set the `LOG_HTTP_HEADERS` environment variable to print request headers to
`STDOUT`. Additionally, set the `LOG_HTTP_BODY` environment variable to print
entire request bodies.

### Server Hostname

Set the `SEND_SERVER_HOSTNAME` environment variable to `false` to prevent the
server from responding with its hostname before echoing the request. The client
may send the `X-Send-Server-Hostname` request header to `true` or `false` to
override this server-wide setting on a per-request basis.

### Arbitrary Headers

Set the `SEND_HEADER_<header-name>` variable to send arbitrary additional
headers in the response. Underscores in the variable name are converted to
hyphens in the header. For example, the following environment variables can be
used to disable CORS:

```bash
SEND_HEADER_ACCESS_CONTROL_ALLOW_ORIGIN="*"
SEND_HEADER_ACCESS_CONTROL_ALLOW_METHODS="*"
SEND_HEADER_ACCESS_CONTROL_ALLOW_HEADERS="*"
```

### WebSocket URL

Set the `WEBSOCKET_ROOT` environment variable to prefix all websocket
requests made by the `.ws` user interface with a specific path.

## Running the server

The examples below show a few different ways of running the server with the HTTP
server bound to a custom TCP port of `10000`.

### Running locally

```
go get -u github.com/jmalloc/echo-server/...
PORT=10000 echo-server
```

```
export GOPROXY=https://goproxy.cn,direct && SEND_SERVER_REMOTE=false PORT=10001 go run cmd/echo-server/main.go
```

## build and package

```bash
go build -o artifacts/echo-server ./cmd/echo-server
# build
GOOS=linux GOARCH=amd64 go build -o artifacts/echo-server_linux_amd64 ./cmd/echo-server
GOOS=linux GOARCH=arm64 go build -o artifacts/echo-server_linux_arm64 ./cmd/echo-server
GOOS=darwin GOARCH=amd64 go build -o artifacts/echo-server_mac_x86_64 ./cmd/echo-server
GOOS=darwin GOARCH=arm64 go build -o artifacts/echo-server_mac_arm64 ./cmd/echo-server
GOOS=windows GOARCH=amd64 go build -o artifacts/echo-server_windows_amd64 ./cmd/echo-server
ls -lah artifacts/*
```

### Running under Docker

To run the latest version as a container:

```
docker run --detach -p 10000:8080 jmalloc/echo-server
```

Or, as a swarm service:

```
docker service create --publish 10000:8080 jmalloc/echo-server
```

The docker container can be built locally with:

```
make docker
```
