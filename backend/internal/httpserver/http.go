package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	notify chan error
	done   chan struct{}
}

func (s *Server) start(ctx context.Context) {
	go func() {
		<-ctx.Done()
		// send to notify not required, because app.Start
		// handle ctx.Done earlier
		err := s.server.Shutdown(context.Background())
		if err != nil {
			log.Println("[HttpServer] closed with error: %v", err)
		}
	}()
	err := s.server.ListenAndServe()
	// ignore ErrServerClose handling, because in
	// this case - app.Start handle ctx.Done earlier
	s.notify <- fmt.Errorf("ListenAndServe: %w", err)
	close(s.notify)
}

type Handler struct {
	Func http.HandlerFunc
	// Pattern() string
	Endpoint string
	Method   string
}
type Handlers []Handler

func New(ctx context.Context, addr string, handlers Handlers) *Server {
	mux := http.NewServeMux()
	for _, h := range handlers {
		p := h.Method + " " + h.Endpoint
		mux.HandleFunc(p, h.Func)
	}
	s := &Server{
		server: &http.Server{
			Addr:        addr,
			Handler:     mux,
			ReadTimeout: 30 * time.Second,
		},
		notify: make(chan error),
		done:   make(chan struct{}),
	}
	s.start(ctx)
	return s
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Done() <-chan struct{} {
	return s.done
}
