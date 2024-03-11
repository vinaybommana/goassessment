package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	timeutils "github.com/vinaybommana/goassessment/utils"
)

var (
	portFlag int
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// App encapsulates the dependencies and behavior of the application.
type App struct {
	SwatchTimeModel *SwatchTimeModel
	WebSocketCtrl   *WebSocketController
}

// NewApp creates a new instance of the App.
func NewApp() *App {
	model := &SwatchTimeModel{}
	controller := &WebSocketController{Model: model}
	return &App{SwatchTimeModel: model, WebSocketCtrl: controller}
}

// Run initializes and runs the application.
func (app *App) Run(port int) {
	// Routes
	http.HandleFunc("/time", app.SwatchTimeHandler)
	http.HandleFunc("/", app.Handler)
	http.HandleFunc("/ws", app.WebSocketCtrl.HandleConnection)
	http.HandleFunc("/timeupdating", app.WebPageHandler)

	// Server initialization
	addr := ":" + strconv.Itoa(port)
	fmt.Printf("running web server on port: %d\n", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// Handlers
func (app *App) Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello world!")
}

func (app *App) SwatchTimeHandler(w http.ResponseWriter, r *http.Request) {
	// validate request
	if r.URL.Path != "/time" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	swatchTime := app.SwatchTimeModel.GetCurrentSwatchTime()
	if swatchTime == "" {
		swatchTime = app.SwatchTimeModel.PrevSwatchTime
	}

	fmt.Fprintf(w, swatchTime)
}

func (app *App) WebPageHandler(w http.ResponseWriter, r *http.Request) {
	// validate request
	if r.URL.Path != "/timeupdating" {
		http.NotFound(w, r)
		return
	}
	// Serve HTML page with WebSocket connection setup
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
	flag.IntVar(&portFlag, "port", 8000, "port with which the web server should run")
	flag.Parse()

	app := NewApp()
	app.Run(portFlag)
}

// SwatchTimeModel represents the model in the MVC pattern.
type SwatchTimeModel struct {
	PrevSwatchTime string
}

// GetCurrentSwatchTime gets the current Swatch Internet Time.
func (m *SwatchTimeModel) GetCurrentSwatchTime() string {
	currentSwatchTime := timeutils.GetCurrentSwatchTime()
	if m.PrevSwatchTime != currentSwatchTime {
		m.PrevSwatchTime = currentSwatchTime
		return currentSwatchTime
	}
	return ""
}

// WebSocketController represents the controller in the MVC pattern.
type WebSocketController struct {
	Model     *SwatchTimeModel
	connMutex sync.Mutex
}

func (c *WebSocketController) HandleConnection(w http.ResponseWriter, r *http.Request) {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	// Upgrading http connection to Upgrader connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Client Connected.")

	for {
		// Get current Swatch Internet Time
		currentSwatchTime := c.Model.GetCurrentSwatchTime()
		if currentSwatchTime != "" {
			err := conn.WriteMessage(websocket.TextMessage, []byte(currentSwatchTime))
			if err != nil {
				log.Println("Failed to write message:", err)
				return
			}
		}

		// sleep for 10 seconds
		time.Sleep(10 * time.Second)
	}
}
