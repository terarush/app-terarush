package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"go-modular/modules/nodes/domain/entity"
	"go-modular/modules/nodes/domain/repository"
	"go-modular/modules/nodes/dto/request"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type NodeService interface {
	CreateNode(userID uint, req *request.CreateNodeRequest) (*entity.Node, error)
	GetNodeByID(id uint, userID uint) (*entity.Node, error)
	GetNodesByUserID(userID uint) ([]entity.Node, error)
	UpdateNode(id uint, userID uint, req *request.UpdateNodeRequest) (*entity.Node, error)
	DeleteNode(id uint, userID uint) error
	StartNode(id uint, userID uint) error
	StopNode(id uint, userID uint) error
	RestartNode(id uint, userID uint) error
	GetAllNodes() ([]entity.Node, error)
}

type nodeServiceImpl struct {
	repo         repository.NodeRepository
	dockerClient *client.Client
}

func NewNodeService(repo repository.NodeRepository) (NodeService, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	return &nodeServiceImpl{
		repo:         repo,
		dockerClient: dockerClient,
	}, nil
}

func (s *nodeServiceImpl) CreateNode(userID uint, req *request.CreateNodeRequest) (*entity.Node, error) {
	ctx := context.Background()

	// Pull image if not exists
	_, _, err := s.dockerClient.ImageInspectWithRaw(ctx, req.Image)
	if err != nil {
		reader, err := s.dockerClient.ImagePull(ctx, req.Image, image.PullOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to pull image: %w", err)
		}
		defer reader.Close()
		// Wait for pull to complete
		_, _ = reader.Read(make([]byte, 1024))
	}

	// Prepare environment variables
	envVars := []string{}
	for k, v := range req.Environment {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}

	// Prepare port bindings
	exposedPorts := nat.PortSet{
		nat.Port(fmt.Sprintf("%d/tcp", req.InternalPort)): struct{}{},
	}

	portBindings := nat.PortMap{
		nat.Port(fmt.Sprintf("%d/tcp", req.InternalPort)): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: strconv.Itoa(req.Port),
			},
		},
	}

	// Prepare resource limits
	resources := container.Resources{
		NanoCPUs: int64(req.CPULimit * 1e9),
		Memory:   req.MemoryLimit * 1024 * 1024, // Convert MB to bytes
	}

	// Prepare restart policy
	restartPolicy := container.RestartPolicy{
		Name: container.RestartPolicyMode(req.RestartPolicy),
	}
	if req.RestartPolicy == "" {
		restartPolicy.Name = container.RestartPolicyUnlessStopped
	}

	// Create container configuration
	containerConfig := &container.Config{
		Image:        req.Image,
		Env:          envVars,
		ExposedPorts: exposedPorts,
	}

	if req.Command != "" {
		containerConfig.Cmd = []string{"/bin/sh", "-c", req.Command}
	}

	hostConfig := &container.HostConfig{
		PortBindings:  portBindings,
		Resources:     resources,
		RestartPolicy: restartPolicy,
	}

	// Create container
	resp, err := s.dockerClient.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	// Convert environment and volumes to JSON
	envJSON, _ := json.Marshal(req.Environment)
	volJSON, _ := json.Marshal(req.Volumes)

	// Set default values
	if req.CPULimit == 0 {
		req.CPULimit = 1.0
	}
	if req.MemoryLimit == 0 {
		req.MemoryLimit = 512
	}
	if req.RestartPolicy == "" {
		req.RestartPolicy = "unless-stopped"
	}

	// Create node entity
	node := &entity.Node{
		UserID:        userID,
		Name:          req.Name,
		Image:         req.Image,
		ContainerID:   resp.ID,
		Status:        entity.NodeStatusCreated,
		Port:          req.Port,
		InternalPort:  req.InternalPort,
		CPULimit:      req.CPULimit,
		MemoryLimit:   req.MemoryLimit,
		Environment:   string(envJSON),
		Volumes:       string(volJSON),
		Command:       req.Command,
		RestartPolicy: req.RestartPolicy,
	}

	// Save to database
	if err := s.repo.Create(node); err != nil {
		// Cleanup: remove container if database save fails
		_ = s.dockerClient.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})
		return nil, fmt.Errorf("failed to save node: %w", err)
	}

	return node, nil
}

func (s *nodeServiceImpl) GetNodeByID(id uint, userID uint) (*entity.Node, error) {
	node, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("node not found: %w", err)
	}

	if node.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to node")
	}

	// Update status from Docker
	_ = s.updateNodeStatus(node)

	return node, nil
}

func (s *nodeServiceImpl) GetNodesByUserID(userID uint) ([]entity.Node, error) {
	nodes, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Update status for all nodes
	for i := range nodes {
		_ = s.updateNodeStatus(&nodes[i])
	}

	return nodes, nil
}

func (s *nodeServiceImpl) UpdateNode(id uint, userID uint, req *request.UpdateNodeRequest) (*entity.Node, error) {
	node, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("node not found: %w", err)
	}

	if node.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to node")
	}

	// Update node fields
	if req.Name != "" {
		node.Name = req.Name
	}
	if req.CPULimit > 0 {
		node.CPULimit = req.CPULimit
	}
	if req.MemoryLimit > 0 {
		node.MemoryLimit = req.MemoryLimit
	}
	if req.Environment != nil {
		envJSON, _ := json.Marshal(req.Environment)
		node.Environment = string(envJSON)
	}
	if req.Volumes != nil {
		volJSON, _ := json.Marshal(req.Volumes)
		node.Volumes = string(volJSON)
	}
	if req.Command != "" {
		node.Command = req.Command
	}
	if req.RestartPolicy != "" {
		node.RestartPolicy = req.RestartPolicy
	}

	if err := s.repo.Update(node); err != nil {
		return nil, fmt.Errorf("failed to update node: %w", err)
	}

	return node, nil
}

func (s *nodeServiceImpl) DeleteNode(id uint, userID uint) error {
	node, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("node not found: %w", err)
	}

	if node.UserID != userID {
		return fmt.Errorf("unauthorized access to node")
	}

	ctx := context.Background()

	// Stop and remove container
	if node.ContainerID != "" {
		_ = s.dockerClient.ContainerStop(ctx, node.ContainerID, container.StopOptions{})
		_ = s.dockerClient.ContainerRemove(ctx, node.ContainerID, container.RemoveOptions{Force: true})
	}

	// Delete from database
	return s.repo.Delete(id)
}

func (s *nodeServiceImpl) StartNode(id uint, userID uint) error {
	node, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("node not found: %w", err)
	}

	if node.UserID != userID {
		return fmt.Errorf("unauthorized access to node")
	}

	ctx := context.Background()

	if err := s.dockerClient.ContainerStart(ctx, node.ContainerID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	node.Status = entity.NodeStatusRunning
	return s.repo.Update(node)
}

func (s *nodeServiceImpl) StopNode(id uint, userID uint) error {
	node, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("node not found: %w", err)
	}

	if node.UserID != userID {
		return fmt.Errorf("unauthorized access to node")
	}

	ctx := context.Background()

	if err := s.dockerClient.ContainerStop(ctx, node.ContainerID, container.StopOptions{}); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	node.Status = entity.NodeStatusStopped
	return s.repo.Update(node)
}

func (s *nodeServiceImpl) RestartNode(id uint, userID uint) error {
	node, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("node not found: %w", err)
	}

	if node.UserID != userID {
		return fmt.Errorf("unauthorized access to node")
	}

	ctx := context.Background()

	if err := s.dockerClient.ContainerRestart(ctx, node.ContainerID, container.StopOptions{}); err != nil {
		return fmt.Errorf("failed to restart container: %w", err)
	}

	node.Status = entity.NodeStatusRunning
	return s.repo.Update(node)
}

func (s *nodeServiceImpl) GetAllNodes() ([]entity.Node, error) {
	return s.repo.FindAll()
}

func (s *nodeServiceImpl) updateNodeStatus(node *entity.Node) error {
	if node.ContainerID == "" {
		return nil
	}

	ctx := context.Background()
	inspect, err := s.dockerClient.ContainerInspect(ctx, node.ContainerID)
	if err != nil {
		node.Status = entity.NodeStatusError
		_ = s.repo.Update(node)
		return err
	}

	if inspect.State.Running {
		node.Status = entity.NodeStatusRunning
	} else {
		node.Status = entity.NodeStatusStopped
	}

	return s.repo.Update(node)
}
