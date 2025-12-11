package usecase

import (
	"context"
	"strings"

	"github.com/allanCordeiro/talent-db/application/domain"
)

type ListTalentUseCase struct {
	TalentGateway domain.TalentGateway
	Ctx           context.Context
}

func NewListTalentUseCase(ctx context.Context, talentGateway domain.TalentGateway) *ListTalentUseCase {
	return &ListTalentUseCase{
		Ctx:           ctx,
		TalentGateway: talentGateway,
	}
}

type ListTalentsInputDTO struct {
	Limit        int
	Cursor       string
	Name         string
	PossibleRole string
	Tags         []string
}

type TalentDTO struct {
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

type ListTalentsOutputDTO struct {
	Talents    []TalentDTO `json:"talents"`
	NextCursor string      `json:"next_cursor,omitempty"`
}

func (uc *ListTalentUseCase) Execute(input ListTalentsInputDTO) (*ListTalentsOutputDTO, error) {
	if input.Limit <= 0 || input.Limit > 50 {
		input.Limit = 50
	}

	talents, nextCursor, err := uc.TalentGateway.GetTalents(uc.Ctx, input.Limit, input.Cursor)
	if err != nil {
		return &ListTalentsOutputDTO{}, err
	}

	var filtered []domain.Talent
	nameFilter := strings.ToLower(strings.TrimSpace(input.Name))
	possibleRoleFilter := strings.ToLower(strings.TrimSpace(input.PossibleRole))

	for _, t := range talents {
		if nameFilter != "" && !strings.Contains(strings.ToLower(t.FullName), nameFilter) {
			continue
		}
		if possibleRoleFilter != "" && !strings.Contains(strings.ToLower(t.PossibleRole), possibleRoleFilter) {
			continue
		}

		filtered = append(filtered, t)
	}

	var talentDTOs []TalentDTO
	for _, t := range filtered {
		talentDTOs = append(talentDTOs, TalentDTO{
			Id:             t.Id.String(),
			ProfileURL:     t.ProfileURL,
			PossibleRole:   t.PossibleRole,
			FullName:       t.FullName,
			Headline:       t.Headline,
			CurrentCompany: t.CurrentCompany,
			CurrentRole:    t.CurrentRole,
			Tags:           t.Tags,
			Notes:          t.Notes,
			CapturedAt:     t.CapturedAt.String(),
		})
	}

	return &ListTalentsOutputDTO{
		Talents:    talentDTOs,
		NextCursor: nextCursor,
	}, nil
}

// buscar por tags. por enquanto não será implementado
// func normalizeTags(tags []string) []string {
// 	normalized := make([]string, 0, len(tags))
// 	for _, tag := range tags {
// 		tt := strings.ToLower(strings.TrimSpace(tag))
// 		if tt != "" {
// 			continue
// 		}
// 		normalized = append(normalized, tt)
// 	}
// 	return normalized
// }
