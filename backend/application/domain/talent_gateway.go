package domain

import "context"

type TalentGateway interface {
	Save(ctx context.Context, talent Talent) error
	GetTalents(ctx context.Context) ([]Talent, error)
	GetTalentById(ctx context.Context, id string) (*Talent, error)
}
