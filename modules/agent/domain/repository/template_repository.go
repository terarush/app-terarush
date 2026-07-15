package repository

import (
	"context"
	"go-modular/modules/agent/domain/entity"
)

type TemplateRepository interface {
	FindAll(ctx context.Context, filter map[string]any) ([]*entity.AgentTemplate, error)
	FindByID(ctx context.Context, id uint) (*entity.AgentTemplate, error)
	FindByName(ctx context.Context, name string) (*entity.AgentTemplate, error)
	FindByCategory(ctx context.Context, category string) ([]*entity.AgentTemplate, error)
	Create(ctx context.Context, tmpl *entity.AgentTemplate) error
	Update(ctx context.Context, tmpl *entity.AgentTemplate) error
	Delete(ctx context.Context, id uint) error
}
