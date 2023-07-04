package core

import (
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"log"
	"net"
)

type RpcServer func(grpc.ServiceRegistrar)

type NewParam struct {
	dig.In

	Option *OptionRpc
	Server []RpcServer `group:"grpc_server"`
}

type OptionRpc struct {
	Addr string
}

type GrpcService struct {
	e      *grpc.Server
	o      *OptionRpc
	Server []RpcServer
}

func NewRpc(o *OptionRpc, p NewParam) Server {
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
