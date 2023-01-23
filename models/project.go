package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          string    `json:"id"`
	Author      string    `json:"author"`
	Url         string    `json:"url"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Project) GetAllProjects() ([]*Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, author, name, url, description, updated_at, created_at from projects`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var projects []*Project
	for rows.Next() {
		var project Project
		err := rows.Scan(
			&project.ID,
			&project.Author,
			&project.Name,
			&project.Url,
			&project.Description,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &project)
	}
	return projects, nil
}

func (p *Project) GetProjectById(id string) (*Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	// fmt.Println("ID", id)
	query := `
		select id, author, url, name, description, created_at, updated_at 
		from projects 
		where id = $1
	`
	var project Project
	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&project.ID,
		&project.Author,
		&project.Url,
		&project.Name,
		&project.Description,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (p *Project) CreateProject(project Project) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	newId := uuid.New()
	query := `
		insert into projects (id, author, url, name, description, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id
	`
	err := db.QueryRowContext(ctx, query,
		newId,
		project.Author,
		project.Url,
		project.Name,
		project.Description,
		time.Now(),
		time.Now(),
	).Scan(&newId)
	if err != nil {
		return "0", err
	}
	return newId.String(), nil
}

func (p *Project) UpdateProject(project Project, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
		update projects set
		author = $1,
		url = $2,
		name = $3,
		description = $4,
		updated_at = $5
		where id = $6
	`

	_, err := db.ExecContext(
		ctx,
		query,
		project.Author,
		project.Url,
		project.Name,
		project.Description,
		time.Now(),
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *Project) DeleteProject(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `delete from projects where id = $1`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
