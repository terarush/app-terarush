package repository

import (
	"context"
	"errors"
	"go-modular/internal/pkg/database"
	"go-modular/modules/agent/domain/entity"
)

var (
	ErrAgentNotFound = errors.New("agent not found")
)

type AgentRepositoryImpl struct{}

func (r AgentRepositoryImpl) FindAll(ctx context.Context, filter map[string]any) ([]*entity.Agent, error) {
	var agents []*entity.Agent
	query := database.DB.WithContext(ctx)
	for k, v := range filter {
		query = query.Where(k, v)
	}
	result := query.Find(&agents)
	if result.Error != nil {
		return nil, result.Error
	}
	return agents, nil
}

func (r AgentRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Agent, error) {
	var agent entity.Agent
	result := database.DB.WithContext(ctx).First(&agent, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &agent, nil
}

func (r AgentRepositoryImpl) FindByName(ctx context.Context, name string) (*entity.Agent, error) {
	var agent entity.Agent
	result := database.DB.WithContext(ctx).Where("name = ?", name).First(&agent)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, ErrAgentNotFound
		}
		return nil, result.Error
	}
	return &agent, nil
}

func (r AgentRepositoryImpl) FindByStatus(ctx context.Context, status string) ([]*entity.Agent, error) {
	var agents []*entity.Agent
	result := database.DB.WithContext(ctx).Where("status = ?", status).Find(&agents)
	if result.Error != nil {
		return nil, result.Error
	}
	return agents, nil
}

func (r AgentRepositoryImpl) Create(ctx context.Context, agent *entity.Agent) error {
	return database.DB.WithContext(ctx).Create(agent).Error
}

func (r AgentRepositoryImpl) Update(ctx context.Context, agent *entity.Agent) error {
	return database.DB.WithContext(ctx).Save(agent).Error
}

func (r AgentRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Agent{}, id).Error
}

func (r AgentRepositoryImpl) Count(ctx context.Context, filter map[string]any) (int64, error) {
	var count int64
	query := database.DB.WithContext(ctx).Model(&entity.Agent{})
	for k, v := range filter {
		query = query.Where(k, v)
	}
	result := query.Count(&count)
	return count, result.Error
}

func NewAgentRepository() AgentRepository {
	return AgentRepositoryImpl{}
}
