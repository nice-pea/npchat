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
			log.Printf("httpserver - Server - closed with error: %v", err)
		}
	}()
	err := s.server.ListenAndServe()
	// ignore ErrServerClose handling, because in
	// this case - app.Start handle ctx.Done earlier
	s.notify <- fmt.Errorf("httpserver - Server - server.ListenAndServe: %w", err)
	close(s.notify)
}

func New(ctx context.Context, addr string) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /registrations", r.)
	mux.HandleFunc("GET /registrations/:id", r.ConfirmRegistration)
	mux.HandleFunc("GET /customers", c.ListCustomers)
	mux.HandleFunc("PATCH /customers/:id", c.UpdateCustomer)

	s := &Server{
		server: &http.Server{
			Addr:        addr,
			Handler:     mux,
			ReadTimeout: 30 * time.Second,
		},
		notify: make(chan error),
		done: make(chan struct{}),
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
