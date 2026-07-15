package repository

import (
	"context"
	"errors"
	"go-modular/internal/pkg/database"
	"go-modular/modules/agent/domain/entity"
	"time"
)

var (
	ErrLogNotFound = errors.New("log not found")
)

type LogRepositoryImpl struct{}

func (r LogRepositoryImpl) FindAll(ctx context.Context, filter map[string]any, page, limit int) ([]*entity.AgentLog, int64, error) {
	var logs []*entity.AgentLog
	var total int64

	query := database.DB.WithContext(ctx).Model(&entity.AgentLog{})
	for k, v := range filter {
		query = query.Where(k, v)
	}
	query.Count(&total)

	offset := (page - 1) * limit
	result := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs)
	return logs, total, result.Error
}

func (r LogRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.AgentLog, error) {
	var log entity.AgentLog
	result := database.DB.WithContext(ctx).First(&log, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &log, nil
}

func (r LogRepositoryImpl) FindByAgentID(ctx context.Context, agentID uint, limit int) ([]*entity.AgentLog, error) {
	var logs []*entity.AgentLog
	result := database.DB.WithContext(ctx).Where("agent_id = ?", agentID).
		Order("created_at DESC").Limit(limit).Find(&logs)
	return logs, result.Error
}

func (r LogRepositoryImpl) FindByTaskID(ctx context.Context, taskID uint, level string) ([]*entity.AgentLog, error) {
	var logs []*entity.AgentLog
	query := database.DB.WithContext(ctx).Where("task_id = ?", taskID)
	if level != "" {
		query = query.Where("level = ?", level)
	}
	result := query.Order("created_at ASC").Find(&logs)
	return logs, result.Error
}

func (r LogRepositoryImpl) FindByLevel(ctx context.Context, level string, limit int) ([]*entity.AgentLog, error) {
	var logs []*entity.AgentLog
	result := database.DB.WithContext(ctx).Where("level = ?", level).
		Order("created_at DESC").Limit(limit).Find(&logs)
	return logs, result.Error
}

func (r LogRepositoryImpl) Create(ctx context.Context, log *entity.AgentLog) error {
	return database.DB.WithContext(ctx).Create(log).Error
}

func (r LogRepositoryImpl) DeleteOlderThan(ctx context.Context, days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	result := database.DB.WithContext(ctx).Where("created_at < ?", cutoff).Delete(&entity.AgentLog{})
	return result.RowsAffected, result.Error
}

func NewLogRepository() LogRepository {
	return LogRepositoryImpl{}
}
