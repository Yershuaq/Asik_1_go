package grpc

import (
	"context"
	"net"

	"github.com/Yershuaq/Asik_1_go/statistics/internal/usecase"
	statpb "github.com/Yershuaq/Asik_1_go/statistics/proto/stats"
	"google.golang.org/grpc"
)

type Server struct {
	uc usecase.Usecase
}

func (s *Server) GetStats(ctx context.Context, req *statpb.StatsRequest) (*statpb.StatsResponse, error) {
	st, err := s.uc.GetStats(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &statpb.StatsResponse{
		TotalOrders: st.TotalOrders,
		TotalItems:  st.TotalItems,
	}, nil
}

func Run(uc usecase.Usecase) error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	statpb.RegisterStatsServiceServer(srv, &Server{uc: uc})
	return srv.Serve(lis)
}
