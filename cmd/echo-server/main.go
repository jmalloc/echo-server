package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Echo server listening on port %s.\n", port)

	err := http.ListenAndServe(
		":"+port,
		h2c.NewHandler(
			http.HandlerFunc(handler),
			&http2.Server{},
		),
	)
	if err != nil {
		panic(err)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool {
		return true
	},
}

func handler(wr http.ResponseWriter, req *http.Request) {
	if os.Getenv("LOG_HTTP_BODY") != "" || os.Getenv("LOG_HTTP_HEADERS") != "" {
		fmt.Printf("--------  %s | %s %s\n", req.RemoteAddr, req.Method, req.URL)
	} else {
		fmt.Printf("%s | %s %s\n", req.RemoteAddr, req.Method, req.URL)
	}

	if os.Getenv("LOG_HTTP_HEADERS") != "" {
		fmt.Printf("Headers\n")
		//Iterate over all header fields
		for k, v := range req.Header {
			fmt.Printf("%q : %q\n", k, v)
		}
	}

	if os.Getenv("LOG_HTTP_BODY") != "" {
		buf := &bytes.Buffer{}
		buf.ReadFrom(req.Body) // nolint:errcheck

		if buf.Len() != 0 {
			w := hex.Dumper(os.Stdout)
			w.Write(buf.Bytes()) // nolint:errcheck
			w.Close()
		}

		// Replace original body with buffered version so it's still sent to the
		// browser.
		req.Body.Close()
		req.Body = ioutil.NopCloser(
			bytes.NewReader(buf.Bytes()),
		)
	}

	if websocket.IsWebSocketUpgrade(req) {
		serveWebSocket(wr, req)
	} else if req.URL.Path == "/.ws" {
		wr.Header().Add("Content-Type", "text/html")
		wr.WriteHeader(200)
		io.WriteString(wr, websocketHTML) // nolint:errcheck
	} else if req.URL.Path == "/.sse" {
		serveSSE(wr, req)
	} else {
		serveHTTP(wr, req)
	}
}

func serveWebSocket(wr http.ResponseWriter, req *http.Request) {
	connection, err := upgrader.Upgrade(wr, req, nil)
	if err != nil {
		fmt.Printf("%s | %s\n", req.RemoteAddr, err)
		return
	}

	defer connection.Close()
	fmt.Printf("%s | upgraded to websocket\n", req.RemoteAddr)

	var message []byte

	host, err := os.Hostname()
	if err == nil {
		message = []byte(fmt.Sprintf("Request served by %s", host))
	} else {
		message = []byte(fmt.Sprintf("Server hostname unknown: %s", err.Error()))
	}

	err = connection.WriteMessage(websocket.TextMessage, message)
	if err == nil {
		var messageType int

		for {
			messageType, message, err = connection.ReadMessage()
			if err != nil {
				break
			}

			if messageType == websocket.TextMessage {
				fmt.Printf("%s | txt | %s\n", req.RemoteAddr, message)
			} else {
				fmt.Printf("%s | bin | %d byte(s)\n", req.RemoteAddr, len(message))
			}

			err = connection.WriteMessage(messageType, message)
			if err != nil {
				break
			}
		}
	}

	if err != nil {
		fmt.Printf("%s | %s\n", req.RemoteAddr, err)
	}
}

func serveHTTP(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Add("Content-Type", "text/plain")
	wr.WriteHeader(200)

	host, err := os.Hostname()
	if err == nil {
		fmt.Fprintf(wr, "Request served by %s\n\n", host)
	} else {
		fmt.Fprintf(wr, "Server hostname unknown: %s\n\n", err.Error())
	}

	writeRequestHeaders(wr, req)
	io.Copy(wr, req.Body) // nolint:errcheck
}

func serveSSE(wr http.ResponseWriter, req *http.Request) {
	if _, ok := wr.(http.Flusher); !ok {
		http.Error(wr, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	var echo strings.Builder
	writeRequestHeaders(&echo, req)
	io.Copy(&echo, req.Body)
	req.Body.Close()

	wr.Header().Set("Content-Type", "text/event-stream")
	wr.Header().Set("Cache-Control", "no-cache")
	wr.Header().Set("Connection", "keep-alive")
	wr.Header().Set("Access-Control-Allow-Origin", "*")

	// Write an event about the server that is serving this request.
	if host, err := os.Hostname(); err == nil {
		sendSSE(wr, req, map[string]string{
			"event":  "server",
			"server": host,
		})
	}

	// Write an event that echoes back the request.
	sendSSE(wr, req, map[string]string{
		"event":   "echo",
		"request": echo.String(),
	})

	// Then send a counter event every second.
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	counter := 0

	for {
		select {
		case <-req.Context().Done():
			return
		case <-ticker.C:
			counter++
		}

		sendSSE(wr, req, map[string]string{
			"event":   "counter",
			"counter": strconv.Itoa(counter),
		})
	}
}

// sendSSE sends a server-sent event and logs it to the console.
func sendSSE(
	wr http.ResponseWriter,
	req *http.Request,
	ev map[string]string,
) {
	for k, v := range ev {
		lines := strings.Split(v, "\n")

		for _, line := range lines {
			fmt.Fprintf(wr, "%s: %s\n", k, line)
			fmt.Printf("%s | sse | %s: %s\n", req.RemoteAddr, k, line)
		}
	}

	fmt.Fprintf(wr, "\n")
	wr.(http.Flusher).Flush()
}

// writeRequestHeaders writes request headers to w.
func writeRequestHeaders(w io.Writer, req *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", req.Proto, req.Method, req.URL)
	fmt.Fprintln(w, "")

	fmt.Fprintf(w, "Host: %s\n", req.Host)
	for key, values := range req.Header {
		for _, value := range values {
			fmt.Fprintf(w, "%s: %s\n", key, value)
		}
	}

	fmt.Fprintln(w, "")
}
