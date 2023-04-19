package service

import (
	"context"

	"github.com/vanyaio/raketa-backend/internal/types"
	proto "github.com/vanyaio/raketa-backend/proto"
)

type IStorage interface {
	CreateUser(ctx context.Context, user *types.User) (*types.User, error)
	CreateTask(ctx context.Context, task *types.Task) (*types.Task, error)
	DeleteTask(ctx context.Context, task *types.Task) error
	AssignWorker(ctx context.Context, req *proto.AssignRequest) (*types.Task, error)
	CloseTask(ctx context.Context, req *proto.CloseRequest) (*types.Task, error)
	GetOpenTasks(ctx context.Context) ([]*proto.Task, error)
}

type BotService struct {
	proto.UnimplementedRaketaServiceServer
	storage IStorage
}

func NewBotService(storage IStorage) *BotService {
	return &BotService{
		storage: storage,
	}
}

func (s *BotService) SignUp(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {

	_, err := s.storage.CreateUser(ctx, &types.User{
		ID: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &proto.RegisterResponse{}, nil
}

func (s *BotService) CreateTask(ctx context.Context, req *proto.CreateRequest) (*proto.CreateResponse, error) {

	t := &types.Task{
		URL:    req.Url,
		Status: &types.Open,
	}

	task, err := s.storage.CreateTask(ctx, t)
	if err != nil {
		return nil, err
	}

	return &proto.CreateResponse{
		Url: task.URL,
	}, nil
}

func (s *BotService) DeleteTask(ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {

	err := s.storage.DeleteTask(ctx, &types.Task{
		URL: req.Url,
	})
	if err != nil {
		return nil, err
	}

	return &proto.DeleteResponse{
		Message: "deleted successfully",
	}, nil
}

func (s *BotService) AssignWorker(ctx context.Context, req *proto.AssignRequest) (*proto.AssignResponse, error) {
	task, err := s.storage.AssignWorker(ctx, req)
	if err != nil {
		return nil, err
	}

	return &proto.AssignResponse{
		Message: task.URL,
	}, nil
}

func (s *BotService) CloseTask(ctx context.Context, req *proto.CloseRequest) (*proto.CloseResponse, error) {
	task, err := s.storage.CloseTask(ctx, req)
	if err != nil {
		return nil, err
	}

	return &proto.CloseResponse{
		Message: *task.Status,
	}, nil
}

func (s *BotService) GetOpenTasks(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	tasks, err := s.storage.GetOpenTasks(ctx)
	if err != nil {
		return nil, err
	}

	return &proto.GetResponse{
		Tasks: tasks,
	}, nil
}
