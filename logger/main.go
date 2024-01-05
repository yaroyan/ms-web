package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yaroyan/ms/logger/constant"
	"github.com/yaroyan/ms/logger/infrastructure/persistence"
	"github.com/yaroyan/ms/logger/interfaces/api/rest/handler"
	"github.com/yaroyan/ms/logger/interfaces/api/rest/router"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	config := constant.GetConfig()

	mongoClient, err := connectToMongo(config.MongoURI)
	if err != nil {
		log.Panic(err)
	}

	router := router.Router{
		Handler: handler.LogHandler{
			Repository: persistence.LogRepository{
				Client: mongoClient,
			},
		},
	}

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// // register the rpc server
	// err = rpc.Register(new(RPCServer))
	// go app.rpcListen()

	// go app.gRPCListen()

	log.Println("Starting service on port: ", config.RestPort)

	// start web server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.RestPort),
		Handler: router.Routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

// func (app *Config) rpcListen() error {
// 	log.Println("Starting RPC server on port: ", rpcPort)
// 	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
// 	if err != nil {
// 		return err
// 	}
// 	defer listen.Close()

// 	for {
// 		rpcConn, err := listen.Accept()
// 		if err != nil {
// 			return err
// 		}
// 		go rpc.ServeConn(rpcConn)
// 	}
// }

func connectToMongo(uri string) (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to mongo.")

	return c, nil
}
