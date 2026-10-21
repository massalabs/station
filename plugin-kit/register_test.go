package register

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestRegisterPluginMissingID(t *testing.T) {
	originalArgs := os.Args
	t.Cleanup(func() { os.Args = originalArgs })

	os.Args = []string{"plugin"}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = listener.Close() })

	err = RegisterPlugin(listener)
	if !errors.Is(err, ErrMissingPluginID) {
		t.Fatalf("expected ErrMissingPluginID, got %v", err)
	}
}

func TestRegisterSuccess(t *testing.T) {
	var received registerBody

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decoding register body: %v", err)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(server.Close)

	addr := &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1234}

	if err := register(server.URL, "plugin-id", addr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.ID != "plugin-id" || received.URL != "http://127.0.0.1:1234" {
		t.Fatalf("unexpected register body: %+v", received)
	}
}

func TestRegisterFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"boom"}`))
	}))
	t.Cleanup(server.Close)

	addr := &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1234}

	err := register(server.URL, "plugin-id", addr)
	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(err.Error(), "HTTP status: 500") || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("unexpected error: %v", err)
	}
}
