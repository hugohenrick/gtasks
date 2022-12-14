package repository

import (
	"errors"
	"time"

	"github.com/hugohenrick/gtasks/database"
	"github.com/hugohenrick/gtasks/models"
	"github.com/hugohenrick/gtasks/utils"
	"gorm.io/gorm"
)

type ITaskRepository interface {
	FindTasks(task models.Task) ([]models.Task, error)
	FindTaskById(id string) (models.Task, error)
	CreateTask(task models.Task) (models.Task, error)
	UpdateTask(id string, task models.Task) (models.Task, error)
	DeleteTask(id string) (int64, error)
	ExecuteTask(id string, task models.Task) (models.Task, error)
}

type TaskRepository struct {
	Database *gorm.DB
}

var TaskRepositoryServices ITaskRepository

func NewTaskRepository() ITaskRepository {
	return &TaskRepository{Database: database.DB}
}

func (t *TaskRepository) FindTasks(task models.Task) ([]models.Task, error) {
	var tasks []models.Task

	result := t.Database.Preload("User").Find(&tasks, task)

	if result.RowsAffected == 0 {
		return []models.Task{}, errors.New("task data not found")
	}

	return tasks, nil
}

func (t *TaskRepository) FindTaskById(id string) (models.Task, error) {
	var task models.Task

	result := t.Database.Preload("User").First(&task, "id = ?", id)

	if result.RowsAffected == 0 {
		return models.Task{}, errors.New("task data not found")
	}

	return task, nil
}

func (t *TaskRepository) CreateTask(task models.Task) (models.Task, error) {
	result := t.Database.Create(&task)

	if result.RowsAffected == 0 {
		return models.Task{}, errors.New("task not created")
	}

	return task, nil
}

func (t *TaskRepository) UpdateTask(id string, task models.Task) (models.Task, error) {
	t.Database.First(&task, id)

	if task.ID == 0 {
		return models.Task{}, errors.New(utils.TaskNotFound)
	}

	result := t.Database.Save(&task)

	if result.RowsAffected == 0 {
		return models.Task{}, errors.New("task not save")
	}

	return task, nil
}

func (t *TaskRepository) DeleteTask(id string) (int64, error) {
	var deletedTask models.Task

	t.Database.First(&deletedTask, id)
	if deletedTask.ID == 0 {
		return 0, errors.New(utils.TaskNotFound)
	}

	result := t.Database.Where("id = ?", id).Delete(&deletedTask)

	if result.RowsAffected == 0 {
		return 0, errors.New("task not deleted")
	}

	return result.RowsAffected, nil
}

func (t *TaskRepository) ExecuteTask(id string, task models.Task) (models.Task, error) {
	t.Database.First(&task, id)

	if task.ID == 0 {
		return models.Task{}, errors.New(utils.TaskNotFound)
	}

	task.Done = true
	timeNow := time.Now()
	task.FinishedAt = &timeNow

	result := t.Database.Save(&task)

	if result.RowsAffected == 0 {
		return models.Task{}, errors.New("task not save")
	}

	return task, nil
}
