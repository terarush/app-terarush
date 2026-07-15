package repository

import (
	"context"
	"go-modular/modules/agent/domain/entity"
)

type ScheduleRepository interface {
	FindAll(ctx context.Context, filter map[string]any) ([]*entity.AgentSchedule, error)
	FindByID(ctx context.Context, id uint) (*entity.AgentSchedule, error)
	FindActive(ctx context.Context) ([]*entity.AgentSchedule, error)
	FindByAgentID(ctx context.Context, agentID uint) ([]*entity.AgentSchedule, error)
	FindDue(ctx context.Context) ([]*entity.AgentSchedule, error)
	Create(ctx context.Context, sched *entity.AgentSchedule) error
	Update(ctx context.Context, sched *entity.AgentSchedule) error
	Delete(ctx context.Context, id uint) error
}
