package usecase

import (
	"context"
	"testing"

	"github.com/allanCordeiro/talent-db/application/domain"
)

type InMemoryTalentGateway struct {
	talents map[string]domain.Talent
}

func NewInMemoryTalentGateway() *InMemoryTalentGateway {
	return &InMemoryTalentGateway{
		talents: make(map[string]domain.Talent),
	}
}

func (g *InMemoryTalentGateway) Save(ctx context.Context, talent domain.Talent) error {
	g.talents[talent.Id.String()] = talent
	return nil
}
func (g *InMemoryTalentGateway) GetTalents(ctx context.Context, limit int, cursor string) ([]domain.Talent, string, error) {
	var talents []domain.Talent
	for _, t := range g.talents {
		talents = append(talents, t)
	}
	return talents, "", nil
}
func (g *InMemoryTalentGateway) GetTalentById(ctx context.Context, id string) (*domain.Talent, error) {
	if talent, exists := g.talents[id]; exists {
		return &talent, nil
	}
	return nil, nil
}

func TestCreateTalentSuccess(t *testing.T) {
	ctx := context.Background()
	gateway := NewInMemoryTalentGateway()
	useCase := NewCreateTalentUseCase(ctx, gateway)

	input := CreateTalentInputDTO{
		ProfileURL:     "https://linkedin.com/in/test",
		PossibleRole:   "Backend Engineer",
		FullName:       "John Doe",
		Headline:       "Senior Developer",
		CurrentCompany: "Tech Corp",
		CurrentRole:    "Lead Engineer",
		Tags:           []string{"golang", "backend"},
		Notes:          "Great candidate",
	}

	output, err := useCase.Execute(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.Id == "" {
		t.Error("expected non-empty ID")
	}

	if _, exists := gateway.talents[output.Id]; !exists {
		t.Error("talent not saved in gateway")
	}
}

func TestCreateTalentSavedData(t *testing.T) {
	ctx := context.Background()
	gateway := NewInMemoryTalentGateway()
	useCase := NewCreateTalentUseCase(ctx, gateway)

	input := CreateTalentInputDTO{
		FullName:       "Jane Smith",
		PossibleRole:   "Frontend Engineer",
		ProfileURL:     "https://linkedin.com/in/jane",
		Headline:       "React Specialist",
		CurrentCompany: "Web Dev Inc",
		CurrentRole:    "Senior Developer",
		Tags:           []string{"react", "typescript"},
		Notes:          "Excellent skills",
	}

	output, err := useCase.Execute(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	saved := gateway.talents[output.Id]
	if saved.FullName != input.FullName {
		t.Errorf("expected FullName %s, got %s", input.FullName, saved.FullName)
	}
}
