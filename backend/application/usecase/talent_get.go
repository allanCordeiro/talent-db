package usecase

import (
	"context"

	"github.com/allanCordeiro/talent-db/application/domain"
)

type GetTalentUseCase struct {
	TalentGateway domain.TalentGateway
	Ctx           context.Context
}

func NewGetTalentUseCase(ctx context.Context, talentGateway domain.TalentGateway) *GetTalentUseCase {
	return &GetTalentUseCase{
		Ctx:           ctx,
		TalentGateway: talentGateway,
	}
}

type GetTalentInputDTO struct {
	Id string
}

type GetTalentOutputDTO struct {
	Id             string   `json:"id"`
	ProfileURL     string   `json:"profile_url"`
	PossibleRole   string   `json:"possible_role"`
	FullName       string   `json:"full_name"`
	Headline       string   `json:"headline"`
	CurrentCompany string   `json:"current_company"`
	CurrentRole    string   `json:"current_role"`
	Tags           []string `json:"tags"`
	Notes          string   `json:"notes"`
	CapturedAt     string   `json:"captured_at"`
}

func (uc *GetTalentUseCase) Execute(input GetTalentInputDTO) (*GetTalentOutputDTO, error) {
	talent, err := uc.TalentGateway.GetTalentById(uc.Ctx, input.Id)
	if err != nil {
		return nil, err
	}
	output := &GetTalentOutputDTO{
		Id:             talent.Id.String(),
		ProfileURL:     talent.ProfileURL,
		PossibleRole:   talent.PossibleRole,
		FullName:       talent.FullName,
		Headline:       talent.Headline,
		CurrentCompany: talent.CurrentCompany,
		CurrentRole:    talent.CurrentRole,
		Tags:           talent.Tags,
		Notes:          talent.Notes,
		CapturedAt:     talent.CapturedAt.String(),
	}
	return output, nil
}
