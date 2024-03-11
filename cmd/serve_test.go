package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
)

func TestSwatchTimeHandler(t *testing.T) {
	tests := []struct {
		method string
		path   string
		status int
	}{
		{"GET", "/time", http.StatusOK},
		{"POST", "/time", http.StatusMethodNotAllowed},
		{"GET", "/unknown", http.StatusNotFound},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(swatchTimeHandler)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != tt.status {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.status)
		}
	}
}

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Hello world!"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestWsConnectionHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(wsConnectionHandler))
	defer server.Close()

	// Convert the server URL to a WebSocket URL
	wsURL := "ws" + server.URL[4:]

	// Create a WebSocket connection to the test server
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Error establishing WebSocket connection: %v", err)
	}
	defer conn.Close()

	// Read the first message from the WebSocket connection
	_, p, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Error reading message from WebSocket: %v", err)
	}

	// Verify that the received message is not empty
	if len(p) == 0 {
		t.Fatal("Received message is empty")
	}
}
