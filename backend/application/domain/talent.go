package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Talent struct {
	Id             uuid.UUID
	ProfileURL     string
	PossibleRole   string
	FullName       string
	Headline       string
	CurrentCompany string
	CurrentRole    string
	Tags           []string
	Notes          string
	CapturedAt     time.Time
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

func BuildTalent(id uuid.UUID, profileUrl string, possibleRole string, fullName string, headline string, currentCompany string,
	currentRole string, tags []string, notes string, capturedAt time.Time) *Talent {
	return &Talent{
		Id:             id,
		ProfileURL:     profileUrl,
		PossibleRole:   possibleRole,
		FullName:       fullName,
		Headline:       headline,
		CurrentCompany: currentCompany,
		CurrentRole:    currentRole,
		Tags:           tags,
		Notes:          notes,
		CapturedAt:     capturedAt,
	}
}
