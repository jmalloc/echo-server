# Echo Server

A very simple HTTP echo server with support for web-sockets.

- Any messages sent from a web-socket client are echoed.
- Visit `/.ws` for a basic UI to connect and send web-socket messages.
- Requests to any other URL will return the request headers and body.
- The `PORT` environment variable sets the server port.
- No TLS support yet :(

To run as a container:

```
docker run --detach -P jmalloc/echo-server
```

To run as a service:

```
docker service create --publish 8080 jmalloc/echo-server
```
