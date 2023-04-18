package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/pgxpool"
	"github.com/vanyaio/raketa-backend/internal/types"
	botpb "github.com/vanyaio/raketa-backend/proto"
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
	ctx = context.TODO()

	u := &types.User{}

	query := `INSERT INTO users (id) VALUES ($1) RETURNING *`
	if err := s.db.QueryRow(ctx, query, user.ID).Scan(&u.ID); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Storage) CreateTask(ctx context.Context, task *types.Task) (*types.Task, error) {
	ctx = context.TODO()

	t := &types.Task{}

	query := `INSERT INTO tasks (url, assigned_id, status) VALUES ($1, NULL, $2) RETURNING *`

	if err := s.db.QueryRow(ctx, query, task.URL, task.Status).Scan(&t.URL, &t.UserID, &t.Status); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Storage) DeleteTask(ctx context.Context, task *types.Task) error {
	ctx = context.TODO()

	query := `DELETE FROM tasks WHERE url = $1`

	_, err := s.db.Exec(ctx, query, task.URL)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) AssignWorker(ctx context.Context, req *botpb.AssignRequest) (*types.Task, error) {
	ctx = context.TODO()

	query := `UPDATE tasks
		SET assigned_id = COALESCE($1, assigned_id)
		WHERE url = $2
		RETURNING *`

	task := &types.Task{}

	if err := s.db.QueryRow(ctx, query, req.UserId, req.Url).Scan(&task.URL, &task.UserID, &task.Status); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Storage) CloseTask(ctx context.Context, req *botpb.CloseRequest) (*types.Task, error) {
	ctx = context.TODO()

	query := `UPDATE tasks
		SET status = 'closed'
		WHERE url = $1
		RETURNING *`

	task := &types.Task{}

	if err := s.db.QueryRow(ctx, query, req.Url).Scan(&task.URL, &task.UserID, &task.Status); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Storage) GetOpenTasks(ctx context.Context) ([]*botpb.Task, error) {
	ctx = context.TODO()

	query := `SELECT * FROM tasks WHERE status = 'open'`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	tasks := []*botpb.Task{}

	for rows.Next() {
		task := &types.Task{}
		if err := rows.Scan(&task.URL, &task.UserID, &task.Status); err != nil {
			return nil, err
		}
		if task.UserID == nil {
			tasks = append(tasks, &botpb.Task{
				Url:    task.URL,
				Status: *task.Status,
			})
		} else {
			tasks = append(tasks, &botpb.Task{
				Url:    task.URL,
				Status: *task.Status,
				UserId: *task.UserID,
			})
		}
	}

	return tasks, nil
}
