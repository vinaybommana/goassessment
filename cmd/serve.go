package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	timeutils "github.com/vinaybommana/goassessment/utils"

	"github.com/gorilla/websocket"
)

var (
	portFlag int
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func swatchTimeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/time" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, timeutils.GetCurrentSwatchTime())
}

func wsConnectionHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrading http connection to Upgrader connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	fmt.Printf("Client Connected.")

	var prevSwatchTime string
	for {
		// Get current Swatch Internet Time
		currentSwatchTime := timeutils.GetCurrentSwatchTime()
		if prevSwatchTime != currentSwatchTime {
			err := conn.WriteMessage(websocket.TextMessage, []byte(currentSwatchTime))
			if err != nil {
				log.Fatal(err)
				return
			}
		}

		prevSwatchTime = currentSwatchTime
		// sleep for 10 seconds
		time.Sleep(10 * time.Second)
	}
}

func webPageConnectionHandler(w http.ResponseWriter, r *http.Request) {
	// Send HTML content with WebSocket connection setup
	htmlContent := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
		    <meta charset="UTF-8">
		    <meta name="viewport" content="width=device-width, initial-scale=1.0">
		    <title>WebSocket Example</title>
		</head>
		<body>
		    <div id="timestamp"></div>

		    <script>
		        const socket = new WebSocket("ws://" + window.location.host + "/ws");

		        socket.onmessage = function(event) {
		            // Update the content of the 'timestamp' div with the received message
		            document.getElementById("timestamp").innerText = "It is currently " + event.data;
		        };

		        socket.onclose = function(event) {
		            console.error("WebSocket closed:", event);
		        };

		        socket.onerror = function(event) {
		            console.error("WebSocket error:", event);
		        };
		    </script>
		</body>
		</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlContent))
}

func main() {
	flag.IntVar(&portFlag, "port", 8000, "port with which the webserver should run")
	flag.Parse()

	http.HandleFunc("/time", swatchTimeHandler)
	http.HandleFunc("/", handler)
	http.HandleFunc("/ws", wsConnectionHandler)
	http.HandleFunc("/timeupdating", webPageConnectionHandler)

	port := strconv.Itoa(portFlag)
	addr := ":" + port
	fmt.Printf("running webserver on port: %d\n", portFlag)
	log.Fatal(http.ListenAndServe(addr, nil))
}
