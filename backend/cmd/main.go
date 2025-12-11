package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firestore_adapter "github.com/allanCordeiro/talent-db/infra/firestore"
	"github.com/allanCordeiro/talent-db/infra/webserver"
)

// @title Talent API
// @version 1.0
// @description API pessoal para biblioteca de talentos.
// @contact.name   Allan Cordeiro
// @contact.email  allan.cordeiro.santos@gmail.com
// @BasePath /
func main() {

	ctx := context.Background()
	fs, err := firestore.NewClient(ctx, "talent-479621")
	if err != nil {
		log.Fatalf("failed to create firestore client: %v", err)
	}

	talentdb := firestore_adapter.NewTalentDB(fs, "talent-479621")
	webserver.Serve(talentdb)

}
