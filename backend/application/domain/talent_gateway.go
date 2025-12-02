package domain

type TalentGateway interface {
	Save(Talent) error
	GetTalents() ([]Talent, error)
	GetTalentById(id string) (*Talent, error)
}
