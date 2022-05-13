package grpc_service

import (
	"github.com/coderzhuang/core/application"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server func(grpc.ServiceRegistrar)

type NewParam struct {
	dig.In

	Option *Option
	Server []Server `group:"grpc_server"`
}

type Option struct {
	Addr string
}

type GrpcService struct {
	e      *grpc.Server
	o      *Option
	Server []Server
}

func New(o *Option, p NewParam) application.Service {
	return &GrpcService{
		e:      grpc.NewServer(),
		o:      o,
		Server: p.Server,
	}
}

func (s *GrpcService) Run() {
	lis, err := net.Listen("tcp", s.o.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}
	for _, server := range s.Server {
		server(s.e)
	}
	log.Printf("server listening at %v", lis.Addr())
	go func() {
		if err := s.e.Serve(lis); err != nil {
			log.Fatalf("GrpcService Start failed. %+v", err)
			return
		}
	}()
}

func (s *GrpcService) Shutdown() {
	s.e.GracefulStop()
}
