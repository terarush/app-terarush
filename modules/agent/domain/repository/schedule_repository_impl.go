package repository

import (
	"context"
	"errors"
	"go-modular/internal/pkg/database"
	"go-modular/modules/agent/domain/entity"
	"time"
)

var (
	ErrScheduleNotFound = errors.New("schedule not found")
)

type ScheduleRepositoryImpl struct{}

func (r ScheduleRepositoryImpl) FindAll(ctx context.Context, filter map[string]any) ([]*entity.AgentSchedule, error) {
	var schedules []*entity.AgentSchedule
	query := database.DB.WithContext(ctx)
	for k, v := range filter {
		query = query.Where(k, v)
	}
	result := query.Order("next_run_at ASC").Find(&schedules)
	return schedules, result.Error
}

func (r ScheduleRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.AgentSchedule, error) {
	var sched entity.AgentSchedule
	result := database.DB.WithContext(ctx).First(&sched, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &sched, nil
}

func (r ScheduleRepositoryImpl) FindActive(ctx context.Context) ([]*entity.AgentSchedule, error) {
	var schedules []*entity.AgentSchedule
	result := database.DB.WithContext(ctx).
		Where("status = ?", "active").
		Where("max_runs = 0 OR run_count < max_runs").
		Find(&schedules)
	return schedules, result.Error
}

func (r ScheduleRepositoryImpl) FindByAgentID(ctx context.Context, agentID uint) ([]*entity.AgentSchedule, error) {
	var schedules []*entity.AgentSchedule
	result := database.DB.WithContext(ctx).Where("agent_id = ?", agentID).Find(&schedules)
	return schedules, result.Error
}

func (r ScheduleRepositoryImpl) FindDue(ctx context.Context) ([]*entity.AgentSchedule, error) {
	var schedules []*entity.AgentSchedule
	now := time.Now()
	result := database.DB.WithContext(ctx).
		Where("status = ?", "active").
		Where("next_run_at IS NULL OR next_run_at <= ?", now).
		Where("max_runs = 0 OR run_count < max_runs").
		Find(&schedules)
	return schedules, result.Error
}

func (r ScheduleRepositoryImpl) Create(ctx context.Context, sched *entity.AgentSchedule) error {
	return database.DB.WithContext(ctx).Create(sched).Error
}

func (r ScheduleRepositoryImpl) Update(ctx context.Context, sched *entity.AgentSchedule) error {
	return database.DB.WithContext(ctx).Save(sched).Error
}

func (r ScheduleRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.AgentSchedule{}, id).Error
}

func NewScheduleRepository() ScheduleRepository {
	return ScheduleRepositoryImpl{}
}
