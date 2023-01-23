package controllers

import (
	"encoding/json"
	"net/http"
	"portfolio-api/helpers"
	"portfolio-api/models"

	"github.com/go-chi/chi/v5"
)

// GET/api/v1/projects
func GetProjects(w http.ResponseWriter, r *http.Request) {
	var projects models.Project
	all, err := projects.GetAllProjects()
	if err != nil {
		h.Errorlog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, all)
}

// GET/api/v1/projects/project/{id}
func GetProjectById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	project, err := mod.Project.GetProjectById(id)
	if err != nil {
		helpers.MessageLogs.Errorlog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, project)
}

// POST/api/v1/projects/project
func CreateProject(w http.ResponseWriter, r *http.Request) {
	var p models.Project
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, p)
	id, err := p.CreateProject(p)
	if err != nil {
		h.Errorlog.Println(err)
		newProject, _ := p.GetProjectById(id)
		helpers.WriteJSON(w, http.StatusOK, newProject)
	}
}

// PUT/api/v1/projects/project/:id
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var p models.Project
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, p)
	err = p.UpdateProject(p, id)
	if err != nil {
		helpers.MessageLogs.Errorlog.Println(err)
		return
	}
}

// DELETE/api/v1/projects/project/:id
// Check -> https://pkg.go.dev/net/http
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := p.DeleteProject(id)

	helpers.WriteJSON(w, http.StatusOK, "deleted project")
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, "Error Deleting Project")
	}
}
