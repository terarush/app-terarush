package repository

import (
	"context"
	"errors"
	"go-modular/internal/pkg/database"
	"go-modular/modules/agent/domain/entity"
)

var (
	ErrCapabilityNotFound = errors.New("capability not found")
)

type CapabilityRepositoryImpl struct{}

func (r CapabilityRepositoryImpl) FindAll(ctx context.Context, agentID uint) ([]*entity.AgentCapability, error) {
	var caps []*entity.AgentCapability
	result := database.DB.WithContext(ctx).Where("agent_id = ?", agentID).Find(&caps)
	return caps, result.Error
}

func (r CapabilityRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.AgentCapability, error) {
	var cap entity.AgentCapability
	result := database.DB.WithContext(ctx).First(&cap, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &cap, nil
}

func (r CapabilityRepositoryImpl) FindByAgentAndCapability(ctx context.Context, agentID uint, capability string) (*entity.AgentCapability, error) {
	var cap entity.AgentCapability
	result := database.DB.WithContext(ctx).Where("agent_id = ? AND capability = ?", agentID, capability).First(&cap)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, ErrCapabilityNotFound
		}
		return nil, result.Error
	}
	return &cap, nil
}

func (r CapabilityRepositoryImpl) Create(ctx context.Context, cap *entity.AgentCapability) error {
	return database.DB.WithContext(ctx).Create(cap).Error
}

func (r CapabilityRepositoryImpl) Update(ctx context.Context, cap *entity.AgentCapability) error {
	return database.DB.WithContext(ctx).Save(cap).Error
}

func (r CapabilityRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.AgentCapability{}, id).Error
}

func NewCapabilityRepository() CapabilityRepository {
	return CapabilityRepositoryImpl{}
}
