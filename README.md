# Echo Server

A very simple HTTP echo server with support for web-sockets.

- Any messages sent from a WebSocket client are echoed.
- Visit `/.ws` for a basic UI to connect and send WebSocket messages.
- Requests to any other URL will return the request headers and body.
- The `PORT` environment variable sets the server port.
- No TLS support yet :(
