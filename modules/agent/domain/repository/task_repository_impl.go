package repository

import (
	"context"
	"errors"
	"go-modular/internal/pkg/database"
	"go-modular/modules/agent/domain/entity"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type TaskRepositoryImpl struct{}

func (r TaskRepositoryImpl) FindAll(ctx context.Context, filter map[string]any, page, limit int) ([]*entity.AgentTask, int64, error) {
	var tasks []*entity.AgentTask
	var total int64

	query := database.DB.WithContext(ctx).Model(&entity.AgentTask{})
	for k, v := range filter {
		query = query.Where(k, v)
	}
	query.Count(&total)

	offset := (page - 1) * limit
	result := query.Order("priority DESC, created_at ASC").Offset(offset).Limit(limit).Find(&tasks)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return tasks, total, nil
}

func (r TaskRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.AgentTask, error) {
	var task entity.AgentTask
	result := database.DB.WithContext(ctx).First(&task, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil
}

func (r TaskRepositoryImpl) FindByAgentID(ctx context.Context, agentID uint, status string) ([]*entity.AgentTask, error) {
	var tasks []*entity.AgentTask
	query := database.DB.WithContext(ctx).Where("assigned_to = ?", agentID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	result := query.Order("priority DESC, created_at ASC").Find(&tasks)
	return tasks, result.Error
}

func (r TaskRepositoryImpl) FindByStatus(ctx context.Context, status string, limit int) ([]*entity.AgentTask, error) {
	var tasks []*entity.AgentTask
	result := database.DB.WithContext(ctx).Where("status = ?", status).
		Order("priority DESC, created_at ASC").
		Limit(limit).Find(&tasks)
	return tasks, result.Error
}

func (r TaskRepositoryImpl) FindPending(ctx context.Context, limit int) ([]*entity.AgentTask, error) {
	var tasks []*entity.AgentTask
	result := database.DB.WithContext(ctx).
		Where("status = ?", "pending").
		Where("assigned_to IS NULL").
		Order("priority DESC, created_at ASC").
		Limit(limit).Find(&tasks)
	return tasks, result.Error
}

func (r TaskRepositoryImpl) FindByParentTaskID(ctx context.Context, parentID uint) ([]*entity.AgentTask, error) {
	var tasks []*entity.AgentTask
	result := database.DB.WithContext(ctx).Where("parent_task_id = ?", parentID).Find(&tasks)
	return tasks, result.Error
}

func (r TaskRepositoryImpl) Create(ctx context.Context, task *entity.AgentTask) error {
	return database.DB.WithContext(ctx).Create(task).Error
}

func (r TaskRepositoryImpl) Update(ctx context.Context, task *entity.AgentTask) error {
	return database.DB.WithContext(ctx).Save(task).Error
}

func (r TaskRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.AgentTask{}, id).Error
}

func NewTaskRepository() TaskRepository {
	return TaskRepositoryImpl{}
}
