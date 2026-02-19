package repository

import "go-modular/modules/nodes/domain/entity"

type NodeRepository interface {
	Create(node *entity.Node) error
	FindByID(id uint) (*entity.Node, error)
	FindByUserID(userID uint) ([]entity.Node, error)
	FindByContainerID(containerID string) (*entity.Node, error)
	Update(node *entity.Node) error
	Delete(id uint) error
	FindAll() ([]entity.Node, error)
}
