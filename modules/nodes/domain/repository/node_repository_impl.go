package repository

import (
	"go-modular/modules/nodes/domain/entity"

	"gorm.io/gorm"
)

type nodeRepositoryImpl struct {
	db *gorm.DB
}

func NewNodeRepository(db *gorm.DB) NodeRepository {
	return &nodeRepositoryImpl{db: db}
}

func (r *nodeRepositoryImpl) Create(node *entity.Node) error {
	return r.db.Create(node).Error
}

func (r *nodeRepositoryImpl) FindByID(id uint) (*entity.Node, error) {
	var node entity.Node
	if err := r.db.First(&node, id).Error; err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *nodeRepositoryImpl) FindByUserID(userID uint) ([]entity.Node, error) {
	var nodes []entity.Node
	if err := r.db.Where("user_id = ?", userID).Find(&nodes).Error; err != nil {
		return nil, err
	}
	return nodes, nil
}

func (r *nodeRepositoryImpl) FindByContainerID(containerID string) (*entity.Node, error) {
	var node entity.Node
	if err := r.db.Where("container_id = ?", containerID).First(&node).Error; err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *nodeRepositoryImpl) Update(node *entity.Node) error {
	return r.db.Save(node).Error
}

func (r *nodeRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entity.Node{}, id).Error
}

func (r *nodeRepositoryImpl) FindAll() ([]entity.Node, error) {
	var nodes []entity.Node
	if err := r.db.Find(&nodes).Error; err != nil {
		return nil, err
	}
	return nodes, nil
}
