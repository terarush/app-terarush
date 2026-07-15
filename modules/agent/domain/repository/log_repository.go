package repository

import (
	"context"
	"go-modular/modules/agent/domain/entity"
)

type LogRepository interface {
	FindAll(ctx context.Context, filter map[string]any, page, limit int) ([]*entity.AgentLog, int64, error)
	FindByID(ctx context.Context, id uint) (*entity.AgentLog, error)
	FindByAgentID(ctx context.Context, agentID uint, limit int) ([]*entity.AgentLog, error)
	FindByTaskID(ctx context.Context, taskID uint, level string) ([]*entity.AgentLog, error)
	FindByLevel(ctx context.Context, level string, limit int) ([]*entity.AgentLog, error)
	Create(ctx context.Context, log *entity.AgentLog) error
	DeleteOlderThan(ctx context.Context, days int) (int64, error)
}
