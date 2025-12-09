package usecase

import (
	"context"

	"github.com/allanCordeiro/talent-db/application/domain"
)

type CreateTalentUseCase struct {
	TalentGateway domain.TalentGateway
	Ctx           context.Context
}

func NewCreateTalentUseCase(ctx context.Context, talentGateway domain.TalentGateway) *CreateTalentUseCase {
	return &CreateTalentUseCase{
		Ctx:           ctx,
		TalentGateway: talentGateway,
	}
}

type CreateTalentInputDTO struct {
	ProfileURL     string   `json:"profile_url"`
	PossibleRole   string   `json:"possible_role"`
	FullName       string   `json:"full_name"`
	Headline       string   `json:"headline"`
	CurrentCompany string   `json:"current_company"`
	CurrentRole    string   `json:"current_role"`
	Tags           []string `json:"tags"`
	Notes          string   `json:"notes"`
}

type CreateTalentOutputDTO struct {
	Id string
}

func (uc *CreateTalentUseCase) Execute(input CreateTalentInputDTO) (*CreateTalentOutputDTO, error) {
	talent, err := domain.Create(
		input.ProfileURL,
		input.PossibleRole,
		input.FullName,
		input.Headline,
		input.CurrentCompany,
		input.CurrentRole,
		input.Tags,
		input.Notes,
	)
	if err != nil {
		return nil, err
	}

	err = uc.TalentGateway.Save(uc.Ctx, *talent)
	if err != nil {
		return nil, err
	}

	output := &CreateTalentOutputDTO{
		Id: talent.Id.String(),
	}
	return output, nil
}
