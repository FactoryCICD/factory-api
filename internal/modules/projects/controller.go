package projects

import (
	"net/http"

	"github.com/FactoryCICD/factory-api/internal/datastore"
	"github.com/FactoryCICD/factory-api/internal/errors"
	"github.com/FactoryCICD/factory-api/pkg/log"
	"github.com/FactoryCICD/factory-api/pkg/types"
)

type ProjectController interface {
	GetProjects() ([]types.Project, error)
	GetProject(id string) (types.Project, error)
	CreateProject(newProject types.CreateProject) (string, error)
	UpdateProject(w http.ResponseWriter, r *http.Request)
}

type projectController struct {
	Datastore datastore.ProjectDatastore
	Logger    log.Logger
}

// NewController creates a new project controller
func NewController(ds datastore.ProjectDatastore, logger log.Logger) ProjectController {
	return &projectController{
		Datastore: ds,
		Logger:    logger,
	}
}

func (c *projectController) GetProjects() ([]types.Project, error) {
	projects, err := c.Datastore.GetProjects()
	if err != nil {
		return nil, errors.NewRequestError(err, errors.InternalServerError, "Error getting projects", c.Logger)
	}

	return projects, nil
}

func (c *projectController) GetProject(id string) (types.Project, error) {
	if id == "" {
		return types.Project{}, errors.NewRequestError(nil, errors.BadRequestError, "ID is required", c.Logger)
	}

	project, err := c.Datastore.GetProject(id)
	if err != nil {
		return types.Project{}, errors.NewRequestError(err, errors.InternalServerError, "Error getting project", c.Logger)
	}

	return project, nil
}

func (c *projectController) CreateProject(newProject types.CreateProject) (string, error) {
	// Validate the project
	if newProject.URL == "" {
		return "", errors.NewRequestError(nil, errors.BadRequestError, "URL is required", c.Logger)
	}

	// Create the project
	project := types.Project{
		URL:           newProject.URL,
		WebhookSecret: newProject.WebhookSecret,
	}

	// Create the project
	projectID, err := c.Datastore.CreateProject(project)
	if err != nil {
		return "", errors.NewRequestError(err, errors.InternalServerError, "Error creating project", c.Logger)
	}

	return projectID, nil
}

func (c *projectController) UpdateProject(w http.ResponseWriter, r *http.Request) {
}
