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

### Building

```
make release
docker build . -t echo
```

### Running

To run as a container:

```
docker run -p 8080:8080 echo:latest --rm
```

or with environment variables:
```
docker run -e SERVER_NAME=[NAME] -e LOG_HTTP_BODY=true [-e etc...] -p 8080:8080 echo:latest --rm
```

You can hit the application in a browser at `http://0.0.0.0:8080/.ws` (provided you didn't change the default port).