package usecase

import "github.com/allanCordeiro/talent-db/application/domain"

type CreateTalentUseCase struct {
	TalentGateway domain.TalentGateway
}

func NewCreateTalentUseCase(talentGateway domain.TalentGateway) *CreateTalentUseCase {
	return &CreateTalentUseCase{
		TalentGateway: talentGateway,
	}
}

type CreateTalentInputDTO struct {
	ProfileURL     string
	PossibleRole   string
	FullName       string
	Headline       string
	CurrentCompany string
	CurrentRole    string
	Tags           []string
	Notes          string
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

	err = uc.TalentGateway.Save(*talent)
	if err != nil {
		return nil, err
	}

	output := &CreateTalentOutputDTO{
		Id: talent.Id.String(),
	}
	return output, nil
}
