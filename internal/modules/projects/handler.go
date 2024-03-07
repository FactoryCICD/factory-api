package projects

import (
	"encoding/json"
	"net/http"

	"github.com/FactoryCICD/factory-api/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ProjectHandler interface {
	GetProjects(w http.ResponseWriter, r *http.Request)
	GetProject(w http.ResponseWriter, r *http.Request)
	CreateProject(w http.ResponseWriter, r *http.Request)
	UpdateProject(w http.ResponseWriter, r *http.Request)
}

type projectHandler struct {
	Controller ProjectController
}

// Routes creates a new project handler
func Routes(controller ProjectController) *chi.Mux {
	r := chi.NewRouter()

	h := &projectHandler{
		Controller: controller,
	}

	r.Get("/", h.GetProjects)
	r.Get("/{id}", h.GetProject)
	r.Post("/", h.CreateProject)
	r.Patch("/{id}", h.UpdateProject)

	return r
}

func (h *projectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
}

func (h *projectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
}

func (h *projectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	// Get the project from the request
	project := types.CreateProject{}
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Create the project
	res, err := h.Controller.CreateProject(project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return the project
	render.JSON(w, r, struct {
		ID string `json:"id"`
	}{
		ID: res,
	})
}

func (h *projectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
}
