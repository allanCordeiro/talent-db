package firestore

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/allanCordeiro/talent-db/application/domain"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
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

func (db *TalentDB) GetTalents(ctx context.Context, limit int, cursor string) ([]domain.Talent, string, error) {
	var talents []domain.Talent
	var lastCapturedAt time.Time
	q := db.fsClient.Collection("talents").OrderBy("captured_at", firestore.Desc)

	if cursor != "" {
		var cursorData time.Time
		if decoded, err := base64.StdEncoding.DecodeString(cursor); err == nil {
			if err := cursorData.UnmarshalText(decoded); err == nil {
				q = q.StartAfter(cursorData)
			}
		}
	}

	q.Limit(limit)
	iter := q.Documents(ctx)
	defer iter.Stop()

	for {
		if len(talents) >= limit {
			break
		}

		doc, err := iter.Next()

		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, "", err
		}

		var talent domain.Talent
		err = doc.DataTo(&talent)
		if err != nil {
			continue
		}
		talents = append(talents, talent)
		lastCapturedAt = talent.CapturedAt
	}

	var nextCursor string
	if !lastCapturedAt.IsZero() && len(talents) == limit {
		encoded, err := lastCapturedAt.MarshalText()
		if err != nil {
			return nil, "", err
		}
		nextCursor = base64.StdEncoding.EncodeToString(encoded)
	}

	return talents, nextCursor, nil
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
