package repos

import (
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/models"
	"go.uber.org/zap"
)

var t = models.Task{}

func (r *repo) CreateTask(task *models.Task) (err error) {
	if err = t.CreateTask(task, r.db); err != nil {
		r.logger.Debug("CreateTask ", zap.Any("error: ", err))
	}

	return nil
}

func (r *repo) GetTaskByID(id uint) (task *models.Task, err error) {
	task, err = t.GetTaskByID(uint64(id), r.db)
	if err != nil {
		r.logger.Debug("GetTaskByID ", zap.Any("error: ", err))
		return nil, err
	}
	return
}

func (r *repo) UpdateTask(task *models.Task) (err error) {
	if err = t.UpdateTask(task, r.db); err != nil {
		r.logger.Debug("UpdateTask ", zap.Any("error: ", err))
	}
	return nil
}

func (r *repo) DeleteTask(id uint) (err error) {
	if err = t.DeleteTask(id, r.db); err != nil {
		r.logger.Debug("DeleteTask ", zap.Any("error: ", err))
	}
	return nil
}

func (r *repo) GetAllTasks() (tasks []*models.Task, err error) {
	tasks, err = t.GetAllTasks(r.db)
	if err != nil {
		r.logger.Debug("GetAllTasks ", zap.Any("error: ", err))
	}
	return
}
