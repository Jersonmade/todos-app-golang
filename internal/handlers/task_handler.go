package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Jersonmade/todos-app-golang/internal/models"
	"github.com/Jersonmade/todos-app-golang/internal/storage"
)

type TaskHandler struct {
	storage storage.CRUDRepository
}

func NewTaskHandler(storage storage.CRUDRepository) *TaskHandler {
	return &TaskHandler{storage: storage}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

func convertIDFromPath(r *http.Request) (int, error) {
	taskIdStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	return strconv.Atoi(taskIdStr)
}

func (th *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createTask models.CreateTaskDTO

	err := json.NewDecoder(r.Body).Decode(&createTask)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	task, err := th.storage.Create(createTask)

	if err != nil {
		if errors.Is(err, storage.InvalidTaskData) {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, task)
}

func (th *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks := th.storage.GetAll()

	respondWithJSON(w, http.StatusOK, tasks)
}

func (th *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	taskId, err := convertIDFromPath(r)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task id")
		return
	}

	task, err := th.storage.Get(taskId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (th *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	taskId, err := convertIDFromPath(r)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task id")
		return
	}

	var updateTask models.UpdateTaskDTO

	err = json.NewDecoder(r.Body).Decode(&updateTask)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	task, err := th.storage.Update(taskId, updateTask)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (th *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	taskId, err := convertIDFromPath(r)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task id")
		return
	}

	err = th.storage.Delete(taskId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": fmt.Sprint("Deleted task with taskId ", taskId),
	})
}
