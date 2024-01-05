package grpc

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net"

// 	"github.com/yaroyan/ms/logger/constant"
// 	ml "github.com/yaroyan/ms/logger/domain/model/log"
// 	"github.com/yaroyan/ms/logger/infrastructure/persistence"
// 	"github.com/yaroyan/ms/logger/logs"
// 	"google.golang.org/grpc"
// )

// type LogServer struct {
// 	logs.UnimplementedLogServiceServer
// 	Repository persistence.LogRepository
// }

// func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
// 	input := req.GetLogEntry()

// 	// write the log
// 	logEntry := ml.Log{
// 		Name: input.Name,
// 		Data: input.Data,
// 	}

// 	err := l.Repository.Insert(logEntry)
// 	if err != nil {
// 		res := &logs.LogResponse{Result: "failed"}
// 		return res, err
// 	}

// 	// return response
// 	res := &logs.LogResponse{Result: "logged."}

// 	return res, nil
// }

// func (l *LogServer) gRPCListen() {
// 	gRpcPort := constant.GetConfig().GRpcPort
// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", gRpcPort))
// 	if err != nil {
// 		log.Fatalf("Failed to liten for gRPC: %v", err)
// 	}

// 	s := grpc.NewServer()

// 	logs.RegisterLogServiceServer(s, &LogServer{Repository: l.Repository})

// 	log.Printf("gRPC Server started on port %d", gRpcPort)

// 	if err := s.Serve(lis); err != nil {
// 		log.Fatalf("Failed to listen for gRPC: %v", err)
// 	}
// }
