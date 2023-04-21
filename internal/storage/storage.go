package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/pgxpool"
	"github.com/vanyaio/raketa-backend/internal/types"
)

type pgxIface interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Ping(context.Context) error
	Begin(context.Context) (pgx.Tx, error)
	Close()
}

type Storage struct {
	db pgxIface
}

func NewStorage(db pgxIface) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateUser(ctx context.Context, user *types.User) error {
	query := `INSERT INTO users (id) VALUES ($1)`
	_, err := s.db.Exec(ctx, query, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CreateTask(ctx context.Context, task *types.Task) error {
	query := `INSERT INTO tasks (url, assigned_id, status) VALUES ($1, NULL, $2)`
	_, err := s.db.Exec(ctx, query, task.Url, task.Status)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteTask(ctx context.Context, task *types.Task) error {
	query := `DELETE FROM tasks WHERE url = $1`
	_, err := s.db.Exec(ctx, query, task.Url)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) AssignUser(ctx context.Context, req *types.AssignUserRequest) error {
	query := `UPDATE tasks
		SET assigned_id = COALESCE($1, assigned_id)
		WHERE url = $2`
	_, err := s.db.Exec(ctx, query, req.UserID, req.Url)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CloseTask(ctx context.Context, req *types.CloseTaskRequest) error {
	query := `UPDATE tasks
		SET status = 'closed'
		WHERE url = $1`
	_, err := s.db.Exec(ctx, query, req.Url)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetOpenTasks(ctx context.Context) ([]*types.Task, error) {
	query := `SELECT * FROM tasks WHERE status = 'open'`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*types.Task{}
	for rows.Next() {
		task := &types.Task{}
		if err := rows.Scan(&task.Url, &task.UserID, &task.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
