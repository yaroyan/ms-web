package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yaroyan/ms/gateway/interfaces/api"
)

const port = "80"

func main() {

	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	log.Printf("Starting Broker service on port %s\n", port)

	r := api.Router{
		Handlers: api.Handlers{
			Client:     &http.Client{},
			Connection: rabbitConn,
		},
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r.Routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// do not continue until rabbit is ready
	for {
		c, err := amqp.Dial(os.Getenv("RABBIT_MQ_URI"))
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ.")
			connection = c
			break
		}

		if 5 < counts {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
