package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Echo server listening on port %s.\n", port)
	err := http.ListenAndServe(":"+port, http.HandlerFunc(handler))
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
	fmt.Printf("%s | %s %s\n", req.RemoteAddr, req.Method, req.URL)
	if websocket.IsWebSocketUpgrade(req) {
		serveWebSocket(wr, req)
	} else if req.URL.Path == "/.ws" {
		wr.Header().Add("Content-Type", "text/html")
		wr.WriteHeader(200)
		io.WriteString(wr, websocketHTML)
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

	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil {
			fmt.Printf("%s | %s\n", req.RemoteAddr, err)
			break
		}

		if messageType == websocket.TextMessage {
			fmt.Printf("%s | txt | %s\n", req.RemoteAddr, message)
		} else {
			fmt.Printf("%s | bin | %d byte(s)\n", req.RemoteAddr, len(message))
		}

		err = connection.WriteMessage(messageType, message)
		if err != nil {
			fmt.Printf("%s | %s\n", req.RemoteAddr, err)
			break
		}
	}
}

func serveHTTP(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Add("Content-Type", "text/plain")
	wr.WriteHeader(200)

	fmt.Fprintf(wr, "%s %s %s\n", req.Proto, req.Method, req.URL)
	fmt.Fprintln(wr, "")
	fmt.Fprintf(wr, "Host: %s\n", req.Host)
	for key, values := range req.Header {
		for _, value := range values {
			fmt.Fprintf(wr, "%s: %s\n", key, value)
		}
	}

	fmt.Fprintln(wr, "")
	io.Copy(wr, req.Body)
}
