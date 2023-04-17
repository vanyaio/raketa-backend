package service

import (
	"context"

	"github.com/vanyaio/raketa-backend/internal/types"
	botpb "github.com/vanyaio/raketa-backend/proto"
)

type IStorage interface {
	CreateUser(ctx context.Context, user *types.User) (*types.User, error)
	CreateTask(ctx context.Context, task *types.Task) (*types.Task, error)
	DeleteTask(ctx context.Context, task *types.Task) error
	AssignWorker(ctx context.Context, req *botpb.AssignRequest) (*types.Task, error)
	CloseTask(ctx context.Context, req *botpb.CloseRequest) (*types.Task, error)
	GetOpenTasks(ctx context.Context) ([]*botpb.Task, error)
}

type BotService struct {
	botpb.UnimplementedBotServiceServer
	storage IStorage
}

func NewBotService(storage IStorage) *BotService {
	return &BotService{
		storage: storage,
	}
}

func (s *BotService) SignUp(ctx context.Context, req *botpb.RegisterRequest) (*botpb.RegisterResponse, error) {

	_, err := s.storage.CreateUser(ctx, &types.User{
		ID: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &botpb.RegisterResponse{}, nil
}

func (s *BotService) CreateTask(ctx context.Context, req *botpb.CreateRequest) (*botpb.CreateResponse, error) {
	
	t := &types.Task{
		URL:    req.Url,
		Status: &types.Open,
	}

	task, err := s.storage.CreateTask(ctx, t)
	if err != nil {
		return nil, err
	}

	return &botpb.CreateResponse{
		Url: task.URL,
	}, nil
}

func (s *BotService) DeleteTask(ctx context.Context, req *botpb.DeleteRequest) (*botpb.DeleteResponse, error) {

	err := s.storage.DeleteTask(ctx, &types.Task{
		URL:   req.Url,
	})
	if err != nil {
		return nil, err
	}

	return &botpb.DeleteResponse{
		Message: "deleted successfully",
	}, nil
}

func (s *BotService) AssignWorker(ctx context.Context, req *botpb.AssignRequest) (*botpb.AssignResponse, error) {
	task, err := s.storage.AssignWorker(ctx, req)
	if err != nil {
		return nil, err
	}

	return &botpb.AssignResponse{
		Message: task.URL,
	}, nil
}

func (s *BotService) CloseTask(ctx context.Context, req *botpb.CloseRequest) (*botpb.CloseResponse, error) {
	task, err := s.storage.CloseTask(ctx, req)
	if err != nil {
		return nil, err
	}

	return &botpb.CloseResponse{
		Message: *task.Status,
	}, nil
}

func (s *BotService) GetOpenTasks(ctx context.Context, req *botpb.GetRequest) (*botpb.GetResponse, error) {
	tasks, err := s.storage.GetOpenTasks(ctx)
	if err != nil {
		return nil, err
	}

	return &botpb.GetResponse{
		Tasks: tasks,
	}, nil
}