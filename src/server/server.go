package server

import (
	"context"
	"environment/dump"
	"environment/logger"
	"idgenerator/app"
	"idgenerator/proto/pbidgenerator"
	"mmapcache/cache"
	"net"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

// Server struct
type Server struct {
	pbidgenerator.IdGeneratorServerServer
	mmapCache   *cache.MMapCache
	mmapCacheCh chan proto.Message
}

// NewServer new
func NewServer() *Server {
	s := &Server{
		mmapCacheCh: make(chan proto.Message, 0x1000),
	}
	return s
}

// Run server
func (s *Server) Run(addr string) error {
	logger.Info("start listen... addr:", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("failed to listen, err:", err)
		return err
	}

	srv := grpc.NewServer()
	pbidgenerator.RegisterIdGeneratorServerServer(srv, s)

	if err := srv.Serve(lis); err != nil {
		logger.Error("failed to serve, err:", err)
	}
	return err
}

// GenSingleMessageId implements proto.
func (s *Server) GenSingleMessageId(ctx context.Context, req *pbidgenerator.SingleMessageId) (*pbidgenerator.SingleMessageIdReply, error) {
	logger.Debug("GenSingleMessageId transid:", req.Transid)
	// 网络事件处理计数器，dump会通过配置将当前服务的网络事件吞吐量提交给监控服务
	dump.NetEventRecvIncr(0)
	defer dump.NetEventRecvDecr(0)

	srvIDs, orderIDs, err := app.GetApp().DB.GenSingleMessageID(req.FromUid, req.ToUid, req.Num)
	if nil != err {
		return nil, err
	}

	return &pbidgenerator.SingleMessageIdReply{
		SrvIds:   srvIDs,
		OrderIds: orderIDs,
	}, nil
}
