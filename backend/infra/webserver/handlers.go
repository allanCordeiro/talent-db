package webserver

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

func (h *Handler) ListTalents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	limitParam := r.URL.Query().Get("limit")
	cursorParam := r.URL.Query().Get("cursor")
	nameParam := r.URL.Query().Get("name")
	possibleRoleParam := r.URL.Query().Get("possible_role")
	tagsParam := r.URL.Query()["tags"]

	uc := usecase.NewListTalentUseCase(r.Context(), h.TalentGateway)
	output, err := uc.Execute(usecase.ListTalentsInputDTO{
		Limit:        parseToInt(limitParam, 50),
		Cursor:       cursorParam,
		Name:         nameParam,
		PossibleRole: possibleRoleParam,
		Tags:         tagsParam,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("application error: " + err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(output)
}

func parseToInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}
