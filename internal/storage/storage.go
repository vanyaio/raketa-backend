package storage

import (
	"context"
	"errors"

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
	query := `INSERT INTO users (id, telegram_username) VALUES ($1, $2)`
	_, err := s.db.Exec(ctx, query, user.ID, user.Username)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CreateTask(ctx context.Context, task *types.Task) error {
	query := `INSERT INTO tasks (url, assigned_id, status, price) VALUES ($1, NULL, $2, 0)`
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
	var id *int64
	if err := s.db.QueryRow(ctx, `SELECT id FROM users WHERE telegram_username = $1`, req.Username).Scan(&id); err != nil {
		if err == pgx.ErrNoRows {
			return pgx.ErrNoRows
		}
		return err
	}
	query := `UPDATE tasks
		SET assigned_id = COALESCE($1, assigned_id)
		WHERE url = $2`
	_, err := s.db.Exec(ctx, query, id, req.Url)
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

func (s *Storage) GetUnassignTasks(ctx context.Context) ([]*types.Task, error) {
	query := `SELECT * FROM tasks WHERE status = 'open' AND assigned_id IS NULL`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*types.Task{}
	for rows.Next() {
		task := &types.Task{}
		if err := rows.Scan(&task.Url, &task.UserID, &task.Status, &task.Price); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Storage) CheckUser(ctx context.Context, user *types.User) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT * FROM users WHERE telegram_username = $1)`
	_ = s.db.QueryRow(ctx, query, user.Username).Scan(&exists)
	if !exists {
		return false, errors.New("User doesn't exists")
	}
	return true, nil
}

func (s *Storage) SetTaskPrice(ctx context.Context, req *types.SetTaskPriceRequest) error {
	query := `UPDATE tasks
			SET price = $1
			WHERE url = $2`
	_, err := s.db.Exec(ctx, query, req.Price, req.Url)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUserStats(ctx context.Context, user *types.User) (int64, error) {
	var tasksCount int64
	query := `SELECT COUNT(*) FROM tasks WHERE assigned_id = $1 AND status = 'closed'`
	if err := s.db.QueryRow(ctx, query, user.ID).Scan(&tasksCount); err != nil {
		if err == pgx.ErrNoRows {
			return 0, pgx.ErrNoRows
		}
		return 0, err
	}
	return tasksCount, nil
}
