package repository

import (
	"context"
	"errors"
	"go-modular/internal/pkg/database"
	"go-modular/modules/agent/domain/entity"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

type SessionRepositoryImpl struct{}

func (r SessionRepositoryImpl) FindAll(ctx context.Context, filter map[string]any) ([]*entity.AgentSession, error) {
	var sessions []*entity.AgentSession
	query := database.DB.WithContext(ctx)
	for k, v := range filter {
		query = query.Where(k, v)
	}
	result := query.Order("started_at DESC").Find(&sessions)
	return sessions, result.Error
}

func (r SessionRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.AgentSession, error) {
	var session entity.AgentSession
	result := database.DB.WithContext(ctx).First(&session, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func (r SessionRepositoryImpl) FindByAgentID(ctx context.Context, agentID uint, status string) ([]*entity.AgentSession, error) {
	var sessions []*entity.AgentSession
	query := database.DB.WithContext(ctx).Where("agent_id = ?", agentID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	result := query.Order("started_at DESC").Find(&sessions)
	return sessions, result.Error
}

func (r SessionRepositoryImpl) FindByUserID(ctx context.Context, userID uint, limit int) ([]*entity.AgentSession, error) {
	var sessions []*entity.AgentSession
	result := database.DB.WithContext(ctx).Where("user_id = ?", userID).
		Order("started_at DESC").Limit(limit).Find(&sessions)
	return sessions, result.Error
}

func (r SessionRepositoryImpl) FindActiveByAgentID(ctx context.Context, agentID uint) (*entity.AgentSession, error) {
	var session entity.AgentSession
	result := database.DB.WithContext(ctx).Where("agent_id = ? AND status = ?", agentID, "active").
		First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func (r SessionRepositoryImpl) Create(ctx context.Context, session *entity.AgentSession) error {
	return database.DB.WithContext(ctx).Create(session).Error
}

func (r SessionRepositoryImpl) Update(ctx context.Context, session *entity.AgentSession) error {
	return database.DB.WithContext(ctx).Save(session).Error
}

func (r SessionRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.AgentSession{}, id).Error
}

func NewSessionRepository() SessionRepository {
	return SessionRepositoryImpl{}
}
