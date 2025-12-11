package webserver

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/allanCordeiro/talent-db/application/domain"
	"github.com/allanCordeiro/talent-db/application/usecase"
	_ "github.com/allanCordeiro/talent-db/docs"
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

// CreateTalent godoc
// @Summary Cria um talento
// @Description Cadastra um talento com os dados enviados no corpo da requisição.
// @Tags talents
// @Accept json
// @Produce json
// @Param talent body usecase.CreateTalentInputDTO true "Dados do talento"
// @Success 201 {object} CreateTalentResponse "Recurso criado"
// @Header 201 {string} Location "URL do talento recém-criado"
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} string "internal error"
// @Router /talent [post]
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

// GetTalent godoc
// @Summary Busca um talento
// @Description Retorna os dados completos de um talento específico.
// @Tags talents
// @Produce json
// @Param id path string true "ID do talento"
// @Success 200 {object} usecase.GetTalentOutputDTO
// @Failure 404 {string} string "talent not found"
// @Failure 500 {string} string "internal error"
// @Router /talent/{id} [get]
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

// ListTalents godoc
// @Summary Lista talentos
// @Description Retorna uma lista de talentos com paginação e filtros em memória.
// @Tags talents
// @Accept json
// @Produce json
// @Param limit query int false "Limite de registros por página"
// @Param cursor query string false "Cursor para próxima página"
// @Param name query string false "Filtro por nome (substring, case-insensitive)"
// @Param possible_role query string false "Filtro por possible role (substring, case-insensitive)"
// @Param tags query []string false "Tags (AND) - múltiplos valores ex: ?tags=go&tags=backend"
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal error"
// @Router /talents [get]
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
