package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/pgxpool"
	"github.com/vanyaio/raketa-backend/internal/types"
)

type PgxIface interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Ping(context.Context) error
	Begin(context.Context) (pgx.Tx, error)
	Close()
}

type Storage struct {
	db PgxIface
}

func NewStorage(db PgxIface) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	u := &types.User{}

	query := `INSERT INTO users (id) VALUES ($1) RETURNING *`
	if err := s.db.QueryRow(ctx, query, user.ID).Scan(&u.ID); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Storage) CreateTask(ctx context.Context, task *types.Task) (*types.Task, error) {
	t := &types.Task{}

	query := `INSERT INTO tasks (url, assigned_id, status) VALUES ($1, NULL, $2) RETURNING *`

	if err := s.db.QueryRow(ctx, query, task.Url, task.Status).Scan(&t.Url, &t.UserID, &t.Status); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Storage) DeleteTask(ctx context.Context, task *types.Task) error {
	query := `DELETE FROM tasks WHERE url = $1`

	_, err := s.db.Exec(ctx, query, task.Url)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) AssignUser(ctx context.Context, req *types.AssignUserRequest) (*types.Task, error) {
	query := `UPDATE tasks
		SET assigned_id = COALESCE($1, assigned_id)
		WHERE url = $2
		RETURNING *`

	task := &types.Task{}

	if err := s.db.QueryRow(ctx, query, req.UserID, req.Url).Scan(&task.Url, &task.UserID, &task.Status); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Storage) CloseTask(ctx context.Context, req *types.CloseTaskRequest) (*types.Task, error) {
	query := `UPDATE tasks
		SET status = 'closed'
		WHERE url = $1
		RETURNING *`

	task := &types.Task{}

	if err := s.db.QueryRow(ctx, query, req.Url).Scan(&task.Url, &task.UserID, &task.Status); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Storage) GetOpenTasks(ctx context.Context) ([]*types.Task, error) {
	query := `SELECT * FROM tasks WHERE status = 'open'`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	tasks := []*types.Task{}

	for rows.Next() {
		task := &types.Task{}
		if err := rows.Scan(&task.Url, &task.UserID, &task.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
