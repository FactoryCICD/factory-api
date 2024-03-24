package datastore

import (
	"context"

	dbcontext "github.com/FactoryCICD/factory-api/pkg/db"
	"github.com/FactoryCICD/factory-api/pkg/types"
)

type ProjectDatastore interface {
	GetProjects() ([]types.Project, error)
	GetProject(id string) (types.Project, error)
	CreateProject(project types.Project) (string, error)
	UpdateProject(id string, project types.Project) (types.Project, error)
}

type projectDatastore struct {
	DB *dbcontext.DB
}

// NewDatastore creates a new project datastore
func NewDatastore(db *dbcontext.DB) ProjectDatastore {
	return &projectDatastore{
		DB: db,
	}
}

func (d *projectDatastore) GetProjects() ([]types.Project, error) {
	projects := []types.Project{}

	rows, err := d.DB.DB().Query(context.TODO(), "SELECT id, url, webhook_secret FROM projects")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		project := types.Project{}
		err := rows.Scan(&project.ID, &project.URL, &project.WebhookSecret)
		if err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	return projects, nil
}

func (d *projectDatastore) GetProject(id string) (types.Project, error) {
	project := types.Project{}

	err := d.DB.DB().QueryRow(context.TODO(), "SELECT id, url, webhook_secret FROM projects WHERE id = $1", id).Scan(&project.ID, &project.URL, &project.WebhookSecret)
	if err != nil {
		return types.Project{}, err
	}

	return project, nil
}

func (d *projectDatastore) CreateProject(project types.Project) (string, error) {
	// Create a transaction
	tx, err := d.DB.DB().Begin(context.TODO())
	if err != nil {
		return "", err
	}

	// Create a new project
	err = d.DB.DB().QueryRow(context.TODO(), "INSERT INTO projects (url, webhook_secret) VALUES ($1, $2) RETURNING id", project.URL, project.WebhookSecret).Scan(&project.ID)
	if err != nil {
		tx.Rollback(context.Background())
		return "", err
	}

	// Commit the transaction
	err = tx.Commit(context.Background())
	if err != nil {
		return "", err
	}

	return project.ID, nil
}

func (d *projectDatastore) UpdateProject(id string, project types.Project) (types.Project, error) {
	return types.Project{}, nil
}
