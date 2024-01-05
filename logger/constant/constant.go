package constant

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	RestPort      int
	RpcPort       int
	GRpcPort      int
	MongoURI      string
	MongoUserName string
	MongoPassword string
}

var singleton *Config
var once sync.Once

const (
	RestPort      = "REST_PORT"
	RpcPort       = "RPC_PORT"
	GRpcPort      = "GRPC_PORT"
	MongoURI      = "MONGO_URI"
	MongoUser     = "MONGO_USER"
	MongoPassword = "MONGO_PASSWORD"
)

func GetConfig() *Config {
	once.Do(func() {
		restPort, err := strconv.Atoi(os.Getenv(RestPort))
		if err != nil {
			slog.Error(fmt.Sprintf("%s is undefined", RestPort), err)
			panic(err)
		}

		rpcPort, err := strconv.Atoi(os.Getenv(RpcPort))
		if err != nil {
			slog.Error(fmt.Sprintf("%s is undefined", RpcPort), err)
			panic(err)
		}

		grpcPort, err := strconv.Atoi(os.Getenv(GRpcPort))
		if err != nil {
			slog.Error(fmt.Sprintf("%s is undefined", GRpcPort), err)
			panic(err)
		}

		mongoURI := os.Getenv(MongoURI)
		if mongoURI == "" {
			slog.Error(fmt.Sprintf("%s is undefined", mongoURI))
		}

		mongoUser := os.Getenv(MongoUser)
		if mongoUser == "" {
			slog.Error(fmt.Sprintf("%s is undefined", mongoUser))
		}

		mongoPassword := os.Getenv(MongoPassword)
		if mongoUser == "" {
			slog.Error(fmt.Sprintf("%s is undefined", mongoPassword))
		}

		singleton = &Config{
			RestPort:      restPort,
			RpcPort:       rpcPort,
			GRpcPort:      grpcPort,
			MongoURI:      mongoURI,
			MongoUserName: mongoUser,
			MongoPassword: mongoPassword,
		}
	})
	return singleton
}
