package main

import (
	"context"
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
		log.Println("handling request from", req.RemoteAddr, req.UserAgent())
		handleResponseHeaders(rw, req)
		handleResponseBody(rw, req)
	})
}

func handleResponseHeaders(rw http.ResponseWriter, req *http.Request) {
	for key, values := range req.Header {
		for _, val := range values {
			rw.Header().Set("X-Ingress-Proxy-Kafka-"+key, val)
		}
	}
	rw.Header().Set("Content-Type", "application/xml")
}

func handleResponseBody(rw http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
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
