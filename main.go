package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
		if _, err := io.Copy(rw, req.Body); err != nil {
			log.Println(err)
		}
	})
}

func handleResponseHeaders(rw http.ResponseWriter, req *http.Request) {
	for key, values := range req.Header {
		for _, val := range values {
			rw.Header().Set("Echo-"+key, val)

			if strings.HasPrefix(key, "X-Ingress-Proxy-Kafka-") {
				rw.Header().Set(key, val)
			}
		}
	}
	rw.Header().Set("Content-Type", "application/xml")
	rw.Header().Set("Http-Version", req.Proto)
	rw.Header().Set("Method", req.Method)
	rw.Header().Set("Full-Path", req.URL.Path)
	rw.Header().Set("Query-String", req.URL.Query().Encode())
	rw.Header().Set("Server-Address", req.Host)
	rw.Header().Set("Client-Address", req.RemoteAddr)
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
