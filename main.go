package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	startServer(ctx)
}

func echoHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		header, err := json.Marshal(req.Header)
		if err != nil {
			log.Println("error marshaling request header:", err)
		}
		log.Println(req.RemoteAddr, string(header))
		handleResponseHeaders(rw, req)
		handleResponseBody(rw, req)
	})
}

func handleResponseHeaders(rw http.ResponseWriter, req *http.Request) {
	for key, values := range req.Header {
		for _, val := range values {
			echoKey := "X-Ingress-Proxy-Kafka-" + key
			if rw.Header().Get(echoKey) == "" {
				rw.Header().Set(echoKey, val)
			} else {
				rw.Header().Add(echoKey, val)
			}
		}
	}
	rw.Header().Set("Content-Type", "application/xml")
}

func handleResponseBody(rw http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("error reading request body:", err)
	}
	if string(body[:3]) == "400" {
		rw.WriteHeader(http.StatusBadRequest)
		body = []byte("400 Bad Request")
	}
	if _, err := rw.Write(body); err != nil {
		log.Println("error writing request body:", err)
	}
}

func startServer(ctx context.Context) {
	h := http.Server{
		Addr:    ":8080",
		Handler: echoHandler(),
	}
	go func() {
		<-ctx.Done()
		if err := h.Shutdown(ctx); err != nil {
			log.Println("error shutting down server:", err)
		}
	}()
	if err := h.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
