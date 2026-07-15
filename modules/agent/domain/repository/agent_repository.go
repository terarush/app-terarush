package repository

import (
	"context"
	"go-modular/modules/agent/domain/entity"
)

type AgentRepository interface {
	FindAll(ctx context.Context, filter map[string]any) ([]*entity.Agent, error)
	FindByID(ctx context.Context, id uint) (*entity.Agent, error)
	FindByName(ctx context.Context, name string) (*entity.Agent, error)
	FindByStatus(ctx context.Context, status string) ([]*entity.Agent, error)
	Create(ctx context.Context, agent *entity.Agent) error
	Update(ctx context.Context, agent *entity.Agent) error
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context, filter map[string]any) (int64, error)
}
