// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package market

import (
	"net"
	"os"
	"sync"

	"github.com/skhal/lab/book/irex/market/service"
	"github.com/skhal/lab/book/irex/pb"
	"google.golang.org/grpc"
)

const defaultSocket = "/tmp/irex_market.sock"

// Server is a the gRPC server with MarkerService running via Unix Domain
// Socket.
type Server struct {
	err        error
	service    *service.Service
	grpcServer *grpc.Server
	socket     string
	once       sync.Once
}

// NewServer creates the server with provides market data.
func NewServer(data *pb.Market) *Server {
	return &Server{
		socket:  defaultSocket,
		service: &service.Service{Quotes: data.GetQuotes()},
	}
}

// Err returns last server error that happened during Serve.
func (s *Server) Err() error { return s.err }

// Serve is a non-blocking call to start the server in a separate Go routine to
// listen and serve request on Unix Domain Socket.
// It returns an error if the server already runs or there is an error to open
// the socket. Use [Server.Err] to access the last server error.
func (s *Server) Serve() (err error) {
	s.once.Do(func() {
		err = s.serve()
	})
	return
}

func (s *Server) serve() error {
	lis, err := net.Listen("unix", s.socket)
	if err != nil {
		return err
	}
	s.grpcServer = grpc.NewServer()
	pb.RegisterMarketServiceServer(s.grpcServer, s.service)
	go func() {
		defer os.Remove(s.socket)
		s.err = s.grpcServer.Serve(lis)
	}()
	return nil
}

// Shutdown gracefully shuts down the server. See [grpc.Server.GracefulStop].
func (s *Server) Shutdown() {
	if s.grpcServer == nil {
		return
	}
	defer s.grpcServer.GracefulStop()
	s.grpcServer = nil
}
