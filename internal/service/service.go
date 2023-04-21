package service

import (
	"context"

	"github.com/vanyaio/raketa-backend/internal/types"
	proto "github.com/vanyaio/raketa-backend/proto"
)

type storage interface {
	CreateUser(ctx context.Context, user *types.User) error
	CreateTask(ctx context.Context, task *types.Task) error
	DeleteTask(ctx context.Context, task *types.Task) error
	AssignUser(ctx context.Context, req *types.AssignUserRequest) error
	CloseTask(ctx context.Context, req *types.CloseTaskRequest) error
	GetOpenTasks(ctx context.Context) ([]*types.Task, error)
}

type Service struct {
	proto.UnimplementedRaketaServiceServer
	storage storage
}

func NewBotService(storage storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) SignUp(ctx context.Context, req *proto.SignUpRequest) (*proto.SignUpResponse, error) {
	if err := s.storage.CreateUser(ctx, &types.User{
		ID: req.Id,
	}); err != nil {
		return nil, err
	}

	return &proto.SignUpResponse{}, nil
}

func (s *Service) CreateTask(ctx context.Context, req *proto.CreateTaskRequest) (*proto.CreateTaskResponse, error) {
	t := &types.Task{
		Url:    req.Url,
		Status: types.Open,
	}

	if err := s.storage.CreateTask(ctx, t); err != nil {
		return nil, err
	}

	return &proto.CreateTaskResponse{}, nil
}

func (s *Service) DeleteTask(ctx context.Context, req *proto.DeleteTaskRequest) (*proto.DeleteTaskResponse, error) {
	if err := s.storage.DeleteTask(ctx, &types.Task{
		Url: req.Url,
	}); err != nil {
		return nil, err
	}

	return &proto.DeleteTaskResponse{}, nil
}

func (s *Service) AssignUser(ctx context.Context, req *proto.AssignUserRequest) (*proto.AssignUserResponse, error) {
	if err := s.storage.AssignUser(ctx, &types.AssignUserRequest{
		Url:    req.Url,
		UserID: &req.UserId,
	}); err != nil {
		return nil, err
	}

	return &proto.AssignUserResponse{}, nil
}

func (s *Service) CloseTask(ctx context.Context, req *proto.CloseTaskRequest) (*proto.CloseTaskResponse, error) {
	if err := s.storage.CloseTask(ctx, &types.CloseTaskRequest{
		Url: req.Url,
	}); err != nil {
		return nil, err
	}

	return &proto.CloseTaskResponse{}, nil
}

func (s *Service) GetOpenTasks(ctx context.Context, req *proto.GetOpenTasksRequest) (*proto.GetOpenTasksResponse, error) {
	tasks, err := s.storage.GetOpenTasks(ctx)
	if err != nil {
		return nil, err
	}

	return &proto.GetOpenTasksResponse{
		Tasks: convertTasksToProto(tasks),
	}, nil
}


func convertTasksToProto(tasks []*types.Task) []*proto.Task {
	tasksProto := []*proto.Task{}

	for _, task := range tasks {
		taskProto := &proto.Task{
			Url:    task.Url,
			Status: converStatusToProto(task.Status),
		}
		if task.UserID != nil {
			taskProto.UserId = *task.UserID
		}
		tasksProto = append(tasksProto, taskProto)
	}

	return tasksProto
}

func converStatusToProto(status types.Status) proto.Task_Status {
	switch status {
	case types.Open:
		return proto.Task_OPEN
	case types.Closed:
		return proto.Task_CLOSED
	case types.Declined:
		return proto.Task_DECLINED
	default:
		return proto.Task_UNKNOWN
	}
}