package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"time"
)

// CachedResponse
type CachedResponse struct {
	StatusCode int
	Body       []byte
	Completed  bool
}

// MemoryStore
type MemoryStore struct {
	mu    sync.Mutex
	items map[string]*CachedResponse
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{items: make(map[string]*CachedResponse)}
}

// Get
func (s *MemoryStore) Get(key string) (*CachedResponse, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cr, ok := s.items[key]
	return cr, ok
}

// StartProcessing
func (s *MemoryStore) StartProcessing(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.items[key]; exists {
		return false
	}
	s.items[key] = &CachedResponse{Completed: false}
	return true
}

// Finish
func (s *MemoryStore) Finish(key string, code int, body []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = &CachedResponse{
		StatusCode: code,
		Body:       body,
		Completed:  true,
	}
}

// IdempotencyMiddleware
func IdempotencyMiddleware(store *MemoryStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			http.Error(w, "Missing Idempotency-Key header", http.StatusBadRequest)
			return
		}

		cached, exists := store.Get(key)
		if exists {
			if cached.Completed {
				for k, v := range cached.Body {
					w.Header().Set(strconv.Itoa(k), string(v))
				}
				w.WriteHeader(cached.StatusCode)
				w.Write(cached.Body)
				return
			} else {

				http.Error(w, "Duplicate request in progress", http.StatusConflict)
				return
			}
		}

		if !store.StartProcessing(key) {
			cached, _ := store.Get(key)
			if cached != nil && cached.Completed {
				w.WriteHeader(cached.StatusCode)
				w.Write(cached.Body)
				return
			}
			http.Error(w, "Duplicate request in progress", http.StatusConflict)
			return
		}

		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		store.Finish(key, rec.Code, rec.Body.Bytes())

		for k, v := range rec.Header() {
			for _, val := range v {
				w.Header().Add(k, val)
			}
		}
		w.WriteHeader(rec.Code)
		w.Write(rec.Body.Bytes())
	})
}

// PaymentHandler
func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing started...")
	time.Sleep(2 * time.Second)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := `{"status":"paid","amount":1000,"transaction_id":"d290f1ee-6c54-4b01-90e6-d701748f0851"}`
	w.Write([]byte(response))
	log.Println("Processing completed.")
}

func main() {
	store := NewMemoryStore()
	mux := http.NewServeMux()
	mux.HandleFunc("/pay", PaymentHandler)
	handler := IdempotencyMiddleware(store, mux)

	idempotencyKey := "unique-key-12345"

	fmt.Println("=== Simulating concurrent duplicate requests ===")
	var wg sync.WaitGroup
	concurrentRequests := 7

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodPost, "/pay", nil)
			req.Header.Set("Idempotency-Key", idempotencyKey)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			log.Printf("Request %d: status=%d body=%s", idx, rec.Code, rec.Body.String())
		}(i)
	}

	wg.Wait()

	fmt.Println("\n=== After completion, sending another request with the same key ===")
	req := httptest.NewRequest(http.MethodPost, "/pay", nil)
	req.Header.Set("Idempotency-Key", idempotencyKey)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	log.Printf("Late request: status=%d body=%s", rec.Code, rec.Body.String())

	fmt.Println("\nDone.")
}
