package repository

import (
	"context"
	"go-modular/modules/agent/domain/entity"
)

type CapabilityRepository interface {
	FindAll(ctx context.Context, agentID uint) ([]*entity.AgentCapability, error)
	FindByID(ctx context.Context, id uint) (*entity.AgentCapability, error)
	FindByAgentAndCapability(ctx context.Context, agentID uint, capability string) (*entity.AgentCapability, error)
	Create(ctx context.Context, cap *entity.AgentCapability) error
	Update(ctx context.Context, cap *entity.AgentCapability) error
	Delete(ctx context.Context, id uint) error
}
