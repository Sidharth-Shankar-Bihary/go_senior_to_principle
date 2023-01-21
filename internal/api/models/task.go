package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Assigned string `json:"assigned"`
	Task     string `json:"task"`
}

func (t *Task) CreateTask(task *Task, db *gorm.DB) (err error) {
	err = db.Create(&task).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) GetTaskByID(taskID uint64, db *gorm.DB) (task *Task, err error) {
	task = &Task{}
	task.ID = uint(taskID)
	if result := db.Find(&task); result.Error != nil {
		return nil, result.Error
	}

	return
}

func (t *Task) UpdateTask(task *Task, db *gorm.DB) (err error) {
	err = db.Model(&task).Updates(task).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) DeleteTask(taskID uint, db *gorm.DB) (err error) {
	task := &Task{}
	task.ID = taskID
	err = db.Delete(&task).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) GetAllTasks(db *gorm.DB) (tasks []*Task, err error) {
	if result := db.Find(&tasks); result.Error != nil {
		return nil, result.Error
	}
	return
}
