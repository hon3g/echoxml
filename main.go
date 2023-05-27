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
		b, _ := json.Marshal(req.Header)
		log.Println(req.RemoteAddr, string(b))
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
		log.Println(err)
	}
	if string(body[:3]) == "400" {
		rw.WriteHeader(http.StatusBadRequest)
		body = []byte("400 Bad Request")
	}
	if _, err := rw.Write(body); err != nil {
		log.Println(err)
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
			log.Println(err)
		}
	}()
	if err := h.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
