package domain

import "context"

type TalentGateway interface {
	Save(ctx context.Context, talent Talent) error
	GetTalents(ctx context.Context, limit int, cursor string) ([]Talent, string, error)
	GetTalentById(ctx context.Context, id string) (*Talent, error)
}
