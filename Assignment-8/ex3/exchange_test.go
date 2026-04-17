package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetRate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/convert" {
			t.Errorf("expected path /convert, got %s", r.URL.Path)
		}
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		if from != "USD" || to != "EUR" {
			t.Errorf("expected from=USD&to=EUR, got from=%s&to=%s", from, to)
		}
		resp := RateResponse{
			Base:   "USD",
			Target: "EUR",
			Rate:   0.85,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	rate, err := service.GetRate("USD", "EUR")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rate != 0.85 {
		t.Errorf("expected rate 0.85, got %f", rate)
	}
}

func TestGetRate_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		resp := RateResponse{ErrorMsg: "invalid currency pair"}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	_, err := service.GetRate("USD", "XYZ")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	expected := "api error: invalid currency pair"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestGetRate_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"base": "USD", "target": "EUR", "rate": 0.85`)) // незакрытая скобка
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	_, err := service.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
}

func TestGetRate_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	service.Client.Timeout = 100 * time.Millisecond
	_, err := service.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestGetRate_InternalServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	_, err := service.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	expected := "unexpected status: 500"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestGetRate_EmptyBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	_, err := service.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
}
