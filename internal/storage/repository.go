package storage

import "github.com/Jersonmade/todos-app-golang/internal/models"

type CRUDRepository interface {
	Create(createTaskDto models.CreateTaskDTO) (models.Task, error)
	Get(taskId int) (models.Task, error)
	GetAll() []models.Task
	Update(taskId int, task models.UpdateTaskDTO) (models.Task, error)
	Delete(taskId int) error
}
