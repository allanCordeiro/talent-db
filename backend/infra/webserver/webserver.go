package webserver

import (
	"log"
	"net/http"
	"os"

	"github.com/allanCordeiro/talent-db/application/domain"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Serve(talentGateway domain.TalentGateway) {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	token := os.Getenv("API_TOKEN")
	if token == "" {
		log.Panicf("API_TOKEN environment variable not set")
	}

	handler := NewHandler(talentGateway, token)
	http.HandleFunc("POST /talent", handler.withAuth(handler.CreateTalent))
	http.HandleFunc("GET /talent/{id}", handler.withAuth(handler.GetTalent))
	http.HandleFunc("GET /talents", handler.withAuth(handler.ListTalents))
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("starting server on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
