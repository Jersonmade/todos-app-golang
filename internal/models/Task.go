package models

type Task struct {
	TaskID      int    `json:"task_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type CreateTaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}
