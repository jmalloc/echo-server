# Echo Server

A very simple HTTP echo server with support for websockets.

## Behavior
- Any messages sent from a websocket client are echoed
- Visit `/.ws` for a basic UI to connect and send websocket messages
- Requests to any other URL will return the request headers and body

## Configuration

- The `PORT` environment variable sets the server port, which defaults to `8080`
- Set the `LOG_HTTP_BODY` environment variable to dump request bodies to `STDOUT`

## Running the server

The examples below show a few different ways of running the server with the HTTP
server bound to a custom TCP port of `10000`.

### Running locally

```
GO111MODULE=off go get -u github.com/jmalloc/echo-server/...
PORT=10000 echo-server
```

### Running under Docker

To run as a container:

```
docker run --detach -p 10000:8080 jmalloc/echo-server
```

To run as a service:

```
docker service create --publish 10000:8080 jmalloc/echo-server
```
