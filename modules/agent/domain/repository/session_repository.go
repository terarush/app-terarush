package repository

import (
	"context"
	"go-modular/modules/agent/domain/entity"
)

type SessionRepository interface {
	FindAll(ctx context.Context, filter map[string]any) ([]*entity.AgentSession, error)
	FindByID(ctx context.Context, id uint) (*entity.AgentSession, error)
	FindByAgentID(ctx context.Context, agentID uint, status string) ([]*entity.AgentSession, error)
	FindByUserID(ctx context.Context, userID uint, limit int) ([]*entity.AgentSession, error)
	FindActiveByAgentID(ctx context.Context, agentID uint) (*entity.AgentSession, error)
	Create(ctx context.Context, session *entity.AgentSession) error
	Update(ctx context.Context, session *entity.AgentSession) error
	Delete(ctx context.Context, id uint) error
}
