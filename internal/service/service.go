package service

import (
	"context"

	"github.com/vanyaio/raketa-backend/internal/types"
	raketapb "github.com/vanyaio/raketa-backend/proto"
)

type IStorage interface {
	CreateUser(ctx context.Context, user *types.User) (*types.User, error)
	CreateTask(ctx context.Context, task *types.Task) (*types.Task, error)
	DeleteTask(ctx context.Context, task *types.Task) error
	AssignWorker(ctx context.Context, req *raketapb.AssignRequest) (*types.Task, error)
	CloseTask(ctx context.Context, req *raketapb.CloseRequest) (*types.Task, error)
	GetOpenTasks(ctx context.Context) ([]*raketapb.Task, error)
}

type BotService struct {
	raketapb.UnimplementedRaketaServiceServer
	storage IStorage
}

func NewBotService(storage IStorage) *BotService {
	return &BotService{
		storage: storage,
	}
}

func (s *BotService) SignUp(ctx context.Context, req *raketapb.RegisterRequest) (*raketapb.RegisterResponse, error) {

	_, err := s.storage.CreateUser(ctx, &types.User{
		ID: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &raketapb.RegisterResponse{}, nil
}

func (s *BotService) CreateTask(ctx context.Context, req *raketapb.CreateRequest) (*raketapb.CreateResponse, error) {

	t := &types.Task{
		URL:    req.Url,
		Status: &types.Open,
	}

	task, err := s.storage.CreateTask(ctx, t)
	if err != nil {
		return nil, err
	}

	return &raketapb.CreateResponse{
		Url: task.URL,
	}, nil
}

func (s *BotService) DeleteTask(ctx context.Context, req *raketapb.DeleteRequest) (*raketapb.DeleteResponse, error) {

	err := s.storage.DeleteTask(ctx, &types.Task{
		URL: req.Url,
	})
	if err != nil {
		return nil, err
	}

	return &raketapb.DeleteResponse{
		Message: "deleted successfully",
	}, nil
}

func (s *BotService) AssignWorker(ctx context.Context, req *raketapb.AssignRequest) (*raketapb.AssignResponse, error) {
	task, err := s.storage.AssignWorker(ctx, req)
	if err != nil {
		return nil, err
	}

	return &raketapb.AssignResponse{
		Message: task.URL,
	}, nil
}

func (s *BotService) CloseTask(ctx context.Context, req *raketapb.CloseRequest) (*raketapb.CloseResponse, error) {
	task, err := s.storage.CloseTask(ctx, req)
	if err != nil {
		return nil, err
	}

	return &raketapb.CloseResponse{
		Message: *task.Status,
	}, nil
}

func (s *BotService) GetOpenTasks(ctx context.Context, req *raketapb.GetRequest) (*raketapb.GetResponse, error) {
	tasks, err := s.storage.GetOpenTasks(ctx)
	if err != nil {
		return nil, err
	}

	return &raketapb.GetResponse{
		Tasks: tasks,
	}, nil
}
