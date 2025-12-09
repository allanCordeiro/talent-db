package firestore

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"github.com/allanCordeiro/talent-db/application/domain"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TalentDB struct {
	fsClient *firestore.Client
	project  string
}

func NewTalentDB(client *firestore.Client, projectID string) *TalentDB {
	return &TalentDB{
		fsClient: client,
		project:  projectID,
	}
}

func (db *TalentDB) Save(ctx context.Context, talent domain.Talent) error {
	t, err := db.GetTalentById(ctx, talent.Id.String())
	if err != nil && status.Code(err) != codes.NotFound {
		return err
	}
	if t != nil {
		_, err := db.fsClient.Collection("talents").Doc(talent.Id.String()).Set(ctx, talent)
		if err != nil {
			return err
		}
	}
	_, err = db.fsClient.Collection("talents").Doc(talent.Id.String()).Create(ctx, talent)
	if err != nil {
		return err
	}
	return nil
}

func (db *TalentDB) GetTalents(ctx context.Context) ([]domain.Talent, error) {
	var talents []domain.Talent
	iter := db.fsClient.Collection("talents").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var talent domain.Talent
		err = doc.DataTo(&talent)
		if err != nil {
			return nil, err
		}
		talents = append(talents, talent)
	}
	return talents, nil
}

func (db *TalentDB) GetTalentById(ctx context.Context, id string) (*domain.Talent, error) {
	doc, err := db.fsClient.Collection("talents").Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, errors.New("talent not found")
	}
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, err
	}
	var talent domain.Talent
	err = doc.DataTo(&talent)
	if err != nil {
		return nil, err
	}

	talent.Id, err = uuid.Parse(doc.Ref.ID)
	if err != nil {
		return nil, err
	}
	return &talent, nil
}
