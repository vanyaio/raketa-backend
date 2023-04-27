package service

import (
	"context"

	"github.com/vanyaio/raketa-backend/internal/types"
	"github.com/vanyaio/raketa-backend/pkg/utils"
	proto "github.com/vanyaio/raketa-backend/proto"
)

const (
	adminRole = "ADMIN_RAKETA"
)

type storage interface {
	CreateUser(ctx context.Context, user *types.User) error
	CreateTask(ctx context.Context, task *types.Task) error
	DeleteTask(ctx context.Context, task *types.Task) error
	AssignUser(ctx context.Context, req *types.AssignUserRequest) error
	CloseTask(ctx context.Context, req *types.CloseTaskRequest) error
	GetUnassignTasks(ctx context.Context) ([]*types.Task, error)
	CheckUser(ctx context.Context, user *types.User) (bool, error)
	SetTaskPrice(ctx context.Context, req *types.SetTaskPriceRequest) error
	GetUserStats(ctx context.Context, user *types.User) (int64, error)
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
		ID:       req.Id,
		Username: req.Username,
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
		Url:      req.Url,
		Username: req.Username,
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

func (s *Service) GetUnassignTasks(ctx context.Context, req *proto.GetUnassignTasksRequest) (*proto.GetUnassignTasksResponse, error) {
	tasks, err := s.storage.GetUnassignTasks(ctx)
	if err != nil {
		return nil, err
	}

	return &proto.GetUnassignTasksResponse{
		Tasks: convertTasksToProto(tasks),
	}, nil
}

func (s *Service) GetUserRole(ctx context.Context, req *proto.GetUserRoleRequest) (*proto.GetUserRoleResponse, error) {
	adminName, err := utils.CheckAdminRole(adminRole)
	if err != nil {
		return nil, err
	}
	ok, err := s.storage.CheckUser(ctx, &types.User{
		Username: req.Username,
	})
	if !ok || err != nil {
		return &proto.GetUserRoleResponse{
			Role: proto.GetUserRoleResponse_UNKNOWN,
		}, err
	}
	if req.Username == adminName {
		return &proto.GetUserRoleResponse{
			Role: proto.GetUserRoleResponse_ADMIN,
		}, nil
	}
	return &proto.GetUserRoleResponse{
		Role: proto.GetUserRoleResponse_REGULAR,
	}, nil
}

func (s *Service) SetTaskPrice(ctx context.Context, req *proto.SetTaskPriceRequest) (*proto.SetTaskPriceResponse, error) {
	if err := s.storage.SetTaskPrice(ctx, &types.SetTaskPriceRequest{
		Url:   req.Url,
		Price: req.Price,
	}); err != nil {
		return nil, err
	}

	return &proto.SetTaskPriceResponse{}, nil
}

func (s *Service) GetUserStats(ctx context.Context, req *proto.GetUserStatsRequest) (*proto.GetUserStatsResponse, error) {
	tasksCount, err := s.storage.GetUserStats(ctx, &types.User{
		ID:       req.UserId,
	})
	if err != nil {
		return &proto.GetUserStatsResponse{}, err
	}
	return &proto.GetUserStatsResponse{
		ClosedTasksCount: tasksCount,
	}, nil
}

func convertTasksToProto(tasks []*types.Task) []*proto.Task {
	tasksProto := []*proto.Task{}

	for _, task := range tasks {
		taskProto := &proto.Task{
			Url:    task.Url,
			Status: converStatusToProto(task.Status),
			Price:  task.Price,
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
