package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/Jersonmade/todos-app-golang/internal/models"
	"github.com/Jersonmade/todos-app-golang/internal/storage"
)

func setupHandler() (storage.CRUDRepository, *TaskHandler) {
	taskStorage := storage.NewMemoryStorage()
	taskHandler := NewTaskHandler(taskStorage)
	return taskStorage, taskHandler
}

func TestCreateSuccess(t *testing.T) {
	_, handler := setupHandler()

	task := models.CreateTaskDTO{
		Title:       "Task1",
		Description: "Description1",
	}

	body, _ := json.Marshal(task)
	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatal("expected status: 201, actual: ", rec.Code)
	}

	var createdTask models.Task

	err := json.NewDecoder(rec.Body).Decode(&createdTask)

	if err != nil {
		t.Fatal("error decoding created task: ", err)
	}

	if createdTask.Title != task.Title ||
		createdTask.Description != task.Description ||
		createdTask.Completed != false {
		t.Error("task fields do not match")
	}
}

func TestCreateValidationError(t *testing.T) {
	_, handler := setupHandler()

	task := models.CreateTaskDTO{
		Title:       "",
		Description: "Description1",
	}

	body, _ := json.Marshal(task)
	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatal("expected status: 400, actual: ", rec.Code)
	}

	expectedJSON := `{"error":"field 'title' is required"}`

	if strings.TrimSpace(rec.Body.String()) != expectedJSON {
		t.Error("expected ", expectedJSON, "actual ", rec.Body.String())
	}
}

func TestGetSuccess(t *testing.T) {
	_, handler := setupHandler()

	task := models.CreateTaskDTO{
		Title:       "Task2",
		Description: "Description2",
	}

	body, _ := json.Marshal(task)
	createReq := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()

	handler.Create(createRec, createReq)

	var createdTask models.Task

	err := json.NewDecoder(createRec.Body).Decode(&createdTask)

	if err != nil {
		t.Fatal("error decoding created task: ", err)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/todos/"+strconv.Itoa(createdTask.TaskID), nil)
	getRec := httptest.NewRecorder()

	handler.Get(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatal("expected status 200, actual ", getRec.Code)
	}

	var receivedTask models.Task

	if err := json.NewDecoder(getRec.Body).Decode(&receivedTask); err != nil {
		t.Fatal("error decoding received task : ", err)
	}

	if receivedTask.TaskID != createdTask.TaskID || receivedTask.Title != createdTask.Title {
		t.Error("task fields do not match")
	}
}

func TestGetNotFound(t *testing.T) {
	_, handler := setupHandler()

	getReq := httptest.NewRequest(http.MethodGet, "/todos/111", nil)
	getRec := httptest.NewRecorder()

	handler.Get(getRec, getReq)

	if getRec.Code != http.StatusNotFound {
		t.Fatal("expected status 404, actual", getRec.Code)
	}

	expectedJSON := `{"error":"task not found"}`

	if strings.TrimSpace(getRec.Body.String()) != expectedJSON {
		t.Error("expected: ", expectedJSON, "actual: ", getRec.Body.String())
	}
}

func TestDeleteSuccess(t *testing.T) {
	_, handler := setupHandler()

	task := models.CreateTaskDTO{
		Title:       "Task3",
		Description: "Description3",
	}

	body, _ := json.Marshal(task)
	createReq := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()

	handler.Create(createRec, createReq)

	var createdTask models.Task

	err := json.NewDecoder(createRec.Body).Decode(&createdTask)

	if err != nil {
		t.Fatal("error decoding created task: ", err)
	}

	delReq := httptest.NewRequest(http.MethodDelete, "/todos/"+strconv.Itoa(createdTask.TaskID), nil)
	delRec := httptest.NewRecorder()

	handler.Delete(delRec, delReq)

	if delRec.Code != http.StatusOK {
		t.Fatal("expected status 200, actual ", delRec.Code)
	}

	expectedJSON := `{"message":"Deleted task with taskId 1"}`

	if strings.TrimSpace(delRec.Body.String()) != expectedJSON {
		t.Error("expected ", expectedJSON, "actual ", delRec.Body.String())
	}

	getReq := httptest.NewRequest(http.MethodGet, "/todos/"+strconv.Itoa(createdTask.TaskID), nil)
	getRec := httptest.NewRecorder()

	handler.Get(getRec, getReq)

	if getRec.Code != http.StatusNotFound {
		t.Error("expected status 404, actual ", getRec.Code)
	}
}

func TestUpdateSuccess(t *testing.T) {
	_, handler := setupHandler()

	task := models.CreateTaskDTO{
		Title:       "Task5",
		Description: "Description5",
	}

	body, _ := json.Marshal(task)
	createReq := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()

	handler.Create(createRec, createReq)

	var createdTask models.Task

	if err := json.NewDecoder(createRec.Body).Decode(&createdTask); err != nil {
		t.Fatal("error decoding created task: ", err)
	}

	update := models.UpdateTaskDTO{
		Title:       "Task5",
		Description: "Description5",
		Completed:   true,
	}

	updateBody, _ := json.Marshal(update)
	updateReq := httptest.NewRequest(http.MethodPut, "/todos/"+strconv.Itoa(createdTask.TaskID), bytes.NewReader(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateRec := httptest.NewRecorder()

	handler.Update(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Fatal("expected status 200, actual", updateRec.Code)
	}

	var updatedTask models.Task

	err := json.NewDecoder(updateRec.Body).Decode(&updatedTask)

	if err != nil {
		t.Fatal("error decoding updated task: ", err)
	}

	if updatedTask.Title != update.Title ||
		updatedTask.Description != update.Description ||
		updatedTask.Completed != update.Completed {
		t.Error("task fields do not match")
	}
}
