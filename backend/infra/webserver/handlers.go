package webserver

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/allanCordeiro/talent-db/application/domain"
	"github.com/allanCordeiro/talent-db/application/usecase"
)

type Handler struct {
	TalentGateway domain.TalentGateway
}

func NewHandler(talentGateway domain.TalentGateway) *Handler {
	return &Handler{
		TalentGateway: talentGateway,
	}
}

type CreateTalentResponse struct {
	Value string `json:"value"`
}

func (h *Handler) CreateTalent(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateTalentInputDTO

	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("data decoder error: " + err.Error())
		return
	}

	uc := usecase.NewCreateTalentUseCase(r.Context(), h.TalentGateway)
	output, err := uc.Execute(usecase.CreateTalentInputDTO{
		ProfileURL:     input.ProfileURL,
		PossibleRole:   input.PossibleRole,
		FullName:       input.FullName,
		Headline:       input.Headline,
		CurrentCompany: input.CurrentCompany,
		CurrentRole:    input.CurrentRole,
		Tags:           input.Tags,
		Notes:          input.Notes,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("application error: " + err.Error())
		return
	}

	response := CreateTalentResponse{
		Value: "/talent/" + output.Id,
	}
	w.Header().Add("Location", response.Value)
	w.WriteHeader(http.StatusCreated)

}

func (h *Handler) GetTalent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	input := usecase.GetTalentInputDTO{
		Id: r.PathValue("id"),
	}
	uc := usecase.NewGetTalentUseCase(r.Context(), h.TalentGateway)
	output, err := uc.Execute(input)
	if err != nil {
		if err.Error() == "talent not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("application error: " + err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(output)

}
