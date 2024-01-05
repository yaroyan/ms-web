package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yaroyan/gms/mail/application/usecase"
	"github.com/yaroyan/gms/mail/interfaces/api/rest/handler"
	"github.com/yaroyan/gms/mail/interfaces/api/rest/router"
)

const webPort = "80"

func main() {

	log.Println("Starting mail service on port: ", webPort)

	r := router.Router{
		MailHandler: &handler.MailHandler{
			Usecase: &usecase.Usecase{},
		},
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: r.Routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
