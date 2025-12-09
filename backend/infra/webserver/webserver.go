package webserver

import (
	"log"
	"net/http"
	"os"

	"github.com/allanCordeiro/talent-db/application/domain"
)

func Serve(talentGateway domain.TalentGateway) {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	handler := NewHandler(talentGateway)
	http.HandleFunc("POST /talent", handler.CreateTalent)
	http.HandleFunc("GET /talent/{id}", handler.GetTalent)

	log.Println("starting server on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
