package internalgrpc

import (
	"fmt"
	"net"

	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/logger"
	"github.com/AndreyNagorskiy/otus-banner-rotation/pb/api/bannerrotation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer *grpc.Server
	logger     logger.Logger
}

func NewServer(logger logger.Logger, bannerRotationHandler pb.BannerRotationServiceServer) *Server {
	s := grpc.NewServer()

	// TODO logging middleware
	pb.RegisterBannerRotationServiceServer(s, bannerRotationHandler)

	reflection.Register(s)

	return &Server{logger: logger, grpcServer: s}
}

func (s *Server) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.logger.Info("GRPC server start on " + addr)

	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
	s.logger.Info("GRPC server stopped")
}
