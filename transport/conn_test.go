package transport

import (
	"SBTP/frame"
	"net"
	"testing"
	"time"
)

func TestConnWriteReadFrame(t *testing.T) {
	serverRaw, clientRaw := net.Pipe()
	defer serverRaw.Close()
	defer clientRaw.Close()

	server := NewConn(serverRaw)
	client := NewConn(clientRaw)

	original := &frame.Frame{
		Version: 1,
		Type:    frame.TypeRequest,
		Status:  uint16(frame.StatusOK),
		Metadata: []frame.Header{
			{Key: "method", Value: "GET"},
			{Key: "path", Value: "/users/42"},
		},
		Payload: []byte(`{"hello":"world"}`),
	}
	original.ApplyChecksum()

	errCh := make(chan error, 1)
	go func() {
		errCh <- client.WriteFrame(original)
	}()

	received, err := server.ReadFrame()
	if err != nil {
		t.Fatalf("ReadFrame failed: %v", err)
	}

	if err := <-errCh; err != nil {
		t.Fatalf("WriteFrame failed: %v", err)
	}

	if received.Version != original.Version {
		t.Errorf("Version mismatch: got %d, want %d", received.Version, original.Version)
	}
	if received.Type != original.Type {
		t.Errorf("Type mismatch: got %v, want %v", received.Type, original.Type)
	}
	if received.Status != original.Status {
		t.Errorf("Status mismatch: got %v, want %v", received.Status, original.Status)
	}
	if string(received.Payload) != string(original.Payload) {
		t.Errorf("Payload mismatch: got %s, want %s", received.Payload, original.Payload)
	}

	method, ok := received.GetHeader("method")
	if !ok || method != "GET" {
		t.Errorf("expected method=GET header, got %q (ok=%v)", method, ok)
	}
}

func TestConnReadTimeout(t *testing.T) {
	serverRaw, clientRaw := net.Pipe()
	defer serverRaw.Close()
	defer clientRaw.Close()

	server := NewConn(serverRaw)
	server.SetReadTimeout(50 * time.Millisecond)

	_ = clientRaw

	start := time.Now()
	_, err := server.ReadFrame()
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if elapsed > 200*time.Millisecond {
		t.Errorf("timeout took too long: %v", elapsed)
	}
}

func TestConnMultipleFramesSequential(t *testing.T) {
	serverRaw, clientRaw := net.Pipe()
	defer serverRaw.Close()
	defer clientRaw.Close()

	server := NewConn(serverRaw)
	client := NewConn(clientRaw)

	frames := []*frame.Frame{
		{Type: frame.TypeRequest, Payload: []byte("first")},
		{Type: frame.TypeRequest, Payload: []byte("second")},
		{Type: frame.TypeRequest, Payload: []byte("third")},
	}

	errCh := make(chan error, 1)
	go func() {
		for _, f := range frames {
			if err := client.WriteFrame(f); err != nil {
				errCh <- err
				return
			}
		}
		errCh <- nil
	}()

	for i, want := range frames {
		got, err := server.ReadFrame()
		if err != nil {
			t.Fatalf("ReadFrame %d failed: %v", i, err)
		}
		if string(got.Payload) != string(want.Payload) {
			t.Errorf("frame %d payload mismatch: got %s, want %s", i, got.Payload, want.Payload)
		}
	}

	if err := <-errCh; err != nil {
		t.Fatalf("WriteFrame goroutine failed: %v", err)
	}
}
