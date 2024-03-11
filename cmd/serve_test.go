package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestSwatchTimeModel_GetCurrentSwatchTime(t *testing.T) {
	model := &SwatchTimeModel{}

	// Initial call should return a non-empty string
	firstCall := model.GetCurrentSwatchTime()
	assert.NotEmpty(t, firstCall)

	// Subsequent calls should return an empty string if the time has not changed
	secondCall := model.GetCurrentSwatchTime()
	assert.Empty(t, secondCall)

	// Simulate a change in time
	model.PrevSwatchTime = ""
	thirdCall := model.GetCurrentSwatchTime()
	assert.NotEmpty(t, thirdCall)
}

func TestWebSocketController_HandleConnection(t *testing.T) {
    _, mockConn, _ := setupWebSocketControllerTest()
	defer mockConn.Close()

	// Capture the messages sent through the WebSocket connection
	var receivedMessages []string

	// Start a goroutine to read messages from the WebSocket connection
	go func() {
		for {
			_, message, err := mockConn.ReadMessage()
			if err != nil {
				return
			}
			receivedMessages = append(receivedMessages, string(message))
		}
	}()

	// Wait for a short duration to allow the WebSocketController to write messages
	time.Sleep(3 * time.Second)

	// Check if at least one message has been received
	assert.NotEmpty(t, receivedMessages)
}

func TestApp_SwatchTimeHandler(t *testing.T) {
	app := NewApp()

	// Create a mock HTTP response writer
	mockResponseWriter := httptest.NewRecorder()

	// Create a mock HTTP request
	mockRequest, err := http.NewRequest("GET", "/time", nil)
	assert.NoError(t, err)

	// Call the SwatchTimeHandler
	app.SwatchTimeHandler(mockResponseWriter, mockRequest)

	// Check the response status code
	assert.Equal(t, http.StatusOK, mockResponseWriter.Code)

	// Check the response body (should not be empty)
	assert.NotEmpty(t, mockResponseWriter.Body.String())
}

func TestApp_WebPageHandler(t *testing.T) {
	app := NewApp()

	// Create a mock HTTP response writer
	mockResponseWriter := httptest.NewRecorder()

	// Create a mock HTTP request
	mockRequest, err := http.NewRequest("GET", "/timeupdating", nil)
	assert.NoError(t, err)

	// Call the WebPageHandler
	app.WebPageHandler(mockResponseWriter, mockRequest)

	// Check the response status code
	assert.Equal(t, http.StatusOK, mockResponseWriter.Code)

	// Check the response body (should not be empty)
	assert.NotEmpty(t, mockResponseWriter.Body.String())
}

func setupWebSocketControllerTest() (*WebSocketController, *websocket.Conn, *httptest.ResponseRecorder) {
	model := &SwatchTimeModel{}
	controller := &WebSocketController{Model: model}

	// Create a test server to handle WebSocket connections
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		// Simulate sending a message after 1 second
		time.Sleep(1 * time.Second)
		err = conn.WriteMessage(websocket.TextMessage, []byte("Mock message"))
		if err != nil {
			http.Error(w, "Failed to write message", http.StatusInternalServerError)
			return
		}
	}))

	// Convert the server URL to a WebSocket URL
	wsURL := "ws" + server.URL[4:]

	// Create a WebSocket connection to the test server
	mockConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		panic(err) // Panic in case of initialization failure
	}

	mockResponseWriter := httptest.NewRecorder()

	return controller, mockConn, mockResponseWriter
}
