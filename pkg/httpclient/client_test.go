package httpclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer srv.Close()

	client := New(srv.URL)
	var result map[string]string
	err := client.Get(context.Background(), "/test", &result)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("unexpected: %v", result)
	}
}

func TestPost(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		json.NewEncoder(w).Encode(body)
	}))
	defer srv.Close()

	client := New(srv.URL)
	var result map[string]string
	err := client.Post(context.Background(), "/test", map[string]string{"a": "1"}, &result)
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	if result["a"] != "1" {
		t.Errorf("unexpected: %v", result)
	}
}
