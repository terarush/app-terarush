package repository

import (
	"context"
	"go-modular/modules/agent/domain/entity"
)

type TaskRepository interface {
	FindAll(ctx context.Context, filter map[string]any, page, limit int) ([]*entity.AgentTask, int64, error)
	FindByID(ctx context.Context, id uint) (*entity.AgentTask, error)
	FindByAgentID(ctx context.Context, agentID uint, status string) ([]*entity.AgentTask, error)
	FindByStatus(ctx context.Context, status string, limit int) ([]*entity.AgentTask, error)
	FindPending(ctx context.Context, limit int) ([]*entity.AgentTask, error)
	FindByParentTaskID(ctx context.Context, parentID uint) ([]*entity.AgentTask, error)
	Create(ctx context.Context, task *entity.AgentTask) error
	Update(ctx context.Context, task *entity.AgentTask) error
	Delete(ctx context.Context, id uint) error
}
