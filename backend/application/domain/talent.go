package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Talent struct {
	Id             uuid.UUID `firestore:"-"`
	ProfileURL     string    `firestore:"profile_url"`
	PossibleRole   string    `firestore:"possible_role"`
	FullName       string    `firestore:"full_name"`
	Headline       string    `firestore:"headline"`
	CurrentCompany string    `firestore:"current_company"`
	CurrentRole    string    `firestore:"current_role"`
	Tags           []string  `firestore:"tags"`
	Notes          string    `firestore:"notes"`
	CapturedAt     time.Time `firestore:"captured_at"`
}

func Create(profileUrl string, possibleRole string, fullName string, headline string, currentCompany string,
	currentRole string, tags []string, notes string) (*Talent, error) {
	talent := &Talent{
		Id:             uuid.New(),
		ProfileURL:     profileUrl,
		PossibleRole:   possibleRole,
		FullName:       fullName,
		Headline:       headline,
		CurrentCompany: currentCompany,
		CurrentRole:    currentRole,
		Tags:           tags,
		Notes:          notes,
		CapturedAt:     time.Now().UTC(),
	}

	err := talent.Validate()
	if err != nil {
		return nil, err
	}
	return talent, nil

}

func (t *Talent) Validate() error {
	if t.ProfileURL == "" {
		return errors.New("url is null")
	}
	if t.PossibleRole == "" {
		return errors.New("role is null")
	}
	if t.FullName == "" {
		return errors.New("name is null")
	}
	if t.Headline == "" {
		return errors.New("headline is null")
	}
	return nil
}
