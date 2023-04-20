package storage

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vanyaio/raketa-backend/internal/types"
	"github.com/vanyaio/raketa-backend/pkg/db"
)

func Test_DB(t *testing.T) {
	ctx := context.Background()
	// db conn
	pool, err := db.NewPool(ctx)
	require.NoError(t, err)
	defer pool.Close()

	storage := NewStorage(pool)

	err = createUser(ctx, t, storage)
	require.NoError(t, err)

	err = createTask(ctx, t, storage)
	require.NoError(t, err)

	err = deleteTask(ctx, t, storage)
	require.NoError(t, err)

	err = assignUser(ctx, t, storage)
	require.NoError(t, err)

	err = closeTask(ctx, t, storage)
	require.NoError(t, err)

	err = getOpenTasks(ctx, t, storage)
	require.NoError(t, err)
}

func createUser(ctx context.Context, t *testing.T, storage *Storage) error {
	u := &types.User{
		ID: int64(randomId()),
	}

	err := storage.CreateUser(ctx, u)
	require.NoError(t, err)

	_, err = storage.db.Exec(ctx, `SELECT * FROM users WHERE id = $1`, u.ID)
	require.NoError(t, err)

	storage.db.Exec(ctx, `TRUNCATE users CASCADE`)

	return err
}

func createTask(ctx context.Context, t *testing.T, storage *Storage) error {
	task := &types.Task{
		Url:    randomURL(),
		Status: types.Open,
	}

	err := storage.CreateTask(ctx, task)
	require.NoError(t, err)

	task = &types.Task{}
	_, err = storage.db.Exec(ctx, `SELECT * FROM tasks WHERE url = $1`, task.Url)
	require.NoError(t, err)

	storage.db.Exec(ctx, `TRUNCATE users CASCADE`)

	return err
}

func deleteTask(ctx context.Context, t *testing.T, storage *Storage) error {
	task := &types.Task{
		Url:    randomURL(),
		Status: types.Open,
	}

	err := storage.CreateTask(ctx, task)
	require.NoError(t, err)

	err = storage.DeleteTask(ctx, task)
	require.NoError(t, err)

	var exists bool
	query := `SELECT NOT EXISTS (SELECT * FROM tasks WHERE url = $1)`
	err = storage.db.QueryRow(ctx, query, task.Url).Scan(&exists)
	if err != nil {
		require.Error(t, err)
	}
	require.NoError(t, err)

	return err
}

func assignUser(ctx context.Context, t *testing.T, storage *Storage) error {
	u := &types.User{
		ID: int64(randomId()),
	}

	err := storage.CreateUser(ctx, u)
	require.NoError(t, err)

	task := &types.Task{
		Url:    randomURL(),
		Status: types.Open,
	}

	err = storage.CreateTask(ctx, task)
	require.NoError(t, err)

	req := &types.AssignUserRequest{
		Url:    task.Url,
		UserID: &u.ID,
	}

	err = storage.AssignUser(ctx, req)
	require.NoError(t, err)

	_, err = storage.db.Exec(ctx, `SELECT * FROM tasks WHERE url = $1 AND assigned_id IS NOT NULL`, task.Url)
	require.NoError(t, err)

	storage.db.Exec(ctx, `TRUNCATE users CASCADE`)

	return err
}

func closeTask(ctx context.Context, t *testing.T, storage *Storage) error {
	task := &types.Task{
		Url:    randomURL(),
		Status: types.Open,
	}

	err := storage.CreateTask(ctx, task)
	require.NoError(t, err)

	req := &types.CloseTaskRequest{
		Url: task.Url,
	}

	err = storage.CloseTask(ctx, req)
	require.NoError(t, err)

	_, err = storage.db.Exec(ctx, `SELECT * FROM tasks WHERE url = $1 AND status = 'closed'`, req.Url)
	require.NoError(t, err)

	storage.db.Exec(ctx, `TRUNCATE users CASCADE`)

	return err
}

func getOpenTasks(ctx context.Context, t *testing.T, storage *Storage) error {
	task1 := &types.Task{
		Url:    randomURL() + fmt.Sprintf("%d", randomId()),
		Status: types.Open,
	}

	err := storage.CreateTask(ctx, task1)
	require.NoError(t, err)

	task2 := &types.Task{
		Url:    randomURL() + fmt.Sprintf("%d2", randomId()),
		Status: types.Open,
	}

	err = storage.CreateTask(ctx, task2)
	require.NoError(t, err)

	tasks, err := storage.GetOpenTasks(ctx)
	require.NoError(t, err)
	require.Contains(t, tasks, task1)
	require.Contains(t, tasks, task2)

	_, err = storage.db.Exec(ctx, `SELECT status, COUNT(*) 
				FROM tasks 
				WHERE status = 'open'
				GROUP BY status 
				HAVING COUNT(*) = 2`)
	require.NoError(t, err)

	storage.db.Exec(ctx, `TRUNCATE users CASCADE`)

	return err
}

func randomURL() string {
	b := make([]byte, 16)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", b)
}

func randomId() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(100000)
}
