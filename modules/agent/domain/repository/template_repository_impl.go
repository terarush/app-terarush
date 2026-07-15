package repository

import (
	"context"
	"errors"
	"go-modular/internal/pkg/database"
	"go-modular/modules/agent/domain/entity"
)

var (
	ErrTemplateNotFound = errors.New("template not found")
)

type TemplateRepositoryImpl struct{}

func (r TemplateRepositoryImpl) FindAll(ctx context.Context, filter map[string]any) ([]*entity.AgentTemplate, error) {
	var templates []*entity.AgentTemplate
	query := database.DB.WithContext(ctx)
	for k, v := range filter {
		query = query.Where(k, v)
	}
	result := query.Order("created_at DESC").Find(&templates)
	return templates, result.Error
}

func (r TemplateRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.AgentTemplate, error) {
	var tmpl entity.AgentTemplate
	result := database.DB.WithContext(ctx).First(&tmpl, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tmpl, nil
}

func (r TemplateRepositoryImpl) FindByName(ctx context.Context, name string) (*entity.AgentTemplate, error) {
	var tmpl entity.AgentTemplate
	result := database.DB.WithContext(ctx).Where("name = ?", name).First(&tmpl)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, ErrTemplateNotFound
		}
		return nil, result.Error
	}
	return &tmpl, nil
}

func (r TemplateRepositoryImpl) FindByCategory(ctx context.Context, category string) ([]*entity.AgentTemplate, error) {
	var templates []*entity.AgentTemplate
	result := database.DB.WithContext(ctx).Where("category = ?", category).Find(&templates)
	return templates, result.Error
}

func (r TemplateRepositoryImpl) Create(ctx context.Context, tmpl *entity.AgentTemplate) error {
	return database.DB.WithContext(ctx).Create(tmpl).Error
}

func (r TemplateRepositoryImpl) Update(ctx context.Context, tmpl *entity.AgentTemplate) error {
	return database.DB.WithContext(ctx).Save(tmpl).Error
}

func (r TemplateRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.AgentTemplate{}, id).Error
}

func NewTemplateRepository() TemplateRepository {
	return TemplateRepositoryImpl{}
}
