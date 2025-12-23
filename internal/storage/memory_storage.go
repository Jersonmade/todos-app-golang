package storage

import (
	"errors"
	"sync"

	"github.com/Jersonmade/todos-app-golang/internal/models"
)

var (
	TaskNotFoundError = errors.New("task not found")
	InvalidTaskData   = errors.New("field 'title' is required")
)

type MemoryStorage struct {
	tasks      map[int]models.Task
	nextTaskId int
	mu         sync.RWMutex
}

func NewMemoryStorage() CRUDRepository {
	return &MemoryStorage{
		tasks:      make(map[int]models.Task),
		nextTaskId: 1,
	}
}

func (ms *MemoryStorage) Create(createTaskDto models.CreateTaskDTO) (models.Task, error) {
	if createTaskDto.Title == "" {
		return models.Task{}, InvalidTaskData
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	task := models.Task{
		TaskID:      ms.nextTaskId,
		Title:       createTaskDto.Title,
		Description: createTaskDto.Description,
		Completed:   false,
	}

	ms.tasks[task.TaskID] = task
	ms.nextTaskId++

	return task, nil
}

func (ms *MemoryStorage) Get(taskId int) (models.Task, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	if task, exists := ms.tasks[taskId]; exists {
		return task, nil
	}

	return models.Task{}, TaskNotFoundError
}

func (ms *MemoryStorage) GetAll() []models.Task {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	tasks := make([]models.Task, 0, len(ms.tasks))

	for _, task := range ms.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

func (ms *MemoryStorage) Update(taskId int, updateTaskDto models.UpdateTaskDTO) (models.Task, error) {
	if updateTaskDto.Title == "" {
		return models.Task{}, InvalidTaskData
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.tasks[taskId]; !exists {
		return models.Task{}, TaskNotFoundError
	}

	task := models.Task{
		TaskID:      taskId,
		Title:       updateTaskDto.Title,
		Description: updateTaskDto.Description,
		Completed:   updateTaskDto.Completed,
	}

	ms.tasks[taskId] = task
	return task, nil
}

func (ms *MemoryStorage) Delete(taskId int) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.tasks[taskId]; !exists {
		return TaskNotFoundError
	}

	delete(ms.tasks, taskId)
	return nil
}
