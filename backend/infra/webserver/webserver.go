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

	handler := NewHandler(talentGateway)
	http.HandleFunc("POST /talent", handler.CreateTalent)
	http.HandleFunc("GET /talent/{id}", handler.GetTalent)
	http.HandleFunc("GET /talents", handler.ListTalents)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("starting server on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
