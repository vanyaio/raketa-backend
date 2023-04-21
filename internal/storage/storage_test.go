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

	pool, err := db.NewPool(ctx)
	require.NoError(t, err)
	defer pool.Close()

	storage := NewStorage(pool)

	createUser(ctx, t, storage)

	createTask(ctx, t, storage)

	deleteTask(ctx, t, storage)

	assignUser(ctx, t, storage)

	closeTask(ctx, t, storage)

	getOpenTasks(ctx, t, storage)
}

func createUser(ctx context.Context, t *testing.T, storage *Storage) {
	u := &types.User{
		ID: int64(randomId()),
	}

	err := storage.CreateUser(ctx, u)
	require.NoError(t, err)

	var exists bool
	query := `SELECT EXISTS (SELECT * FROM users WHERE id = $1)`
	err = storage.db.QueryRow(ctx, query, u.ID).Scan(&exists)
	if err != nil {
		require.Error(t, err)
	}
	require.True(t, exists)
	require.NoError(t, err)

	storage.db.Exec(ctx, `TRUNCATE users CASCADE`)
}

func createTask(ctx context.Context, t *testing.T, storage *Storage) {
	task := &types.Task{
		Url:    randomURL(),
		Status: types.Open,
	}

	err := storage.CreateTask(ctx, task)
	require.NoError(t, err)

	var exists bool
	query := `SELECT EXISTS (SELECT * FROM tasks WHERE url = $1)`
	err = storage.db.QueryRow(ctx, query, task.Url).Scan(&exists)
	if err != nil {
		require.Error(t, err)
	}
	require.True(t, exists)
	require.NoError(t, err)

	storage.db.Exec(ctx, `TRUNCATE users CASCADE`)
}

func deleteTask(ctx context.Context, t *testing.T, storage *Storage) {
	task := &types.Task{
		Url:    randomURL(),
		Status: types.Open,
	}

	err := storage.CreateTask(ctx, task)
	require.NoError(t, err)

	err = storage.DeleteTask(ctx, task)
	require.NoError(t, err)

	var notExists bool
	query := `SELECT NOT EXISTS (SELECT * FROM tasks WHERE url = $1)`
	err = storage.db.QueryRow(ctx, query, task.Url).Scan(&notExists)
	if err != nil {
		require.Error(t, err)
	}
	require.True(t, notExists)
	require.NoError(t, err)
}

func assignUser(ctx context.Context, t *testing.T, storage *Storage) {
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

	var exists bool
	query := `SELECT EXISTS (SELECT * FROM tasks WHERE url = $1 AND assigned_id IS NOT NULL)`
	err = storage.db.QueryRow(ctx, query, task.Url).Scan(&exists)
	if err != nil {
		require.Error(t, err)
	}
	require.True(t, exists)
	require.NoError(t, err)

	storage.db.Exec(ctx, `TRUNCATE users CASCADE`)
}

func closeTask(ctx context.Context, t *testing.T, storage *Storage) {
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

	var exists bool
	query := `SELECT EXISTS (SELECT * FROM tasks WHERE url = $1 AND status = 'closed')`
	err = storage.db.QueryRow(ctx, query, req.Url).Scan(&exists)
	if err != nil {
		require.Error(t, err)
	}
	require.True(t, exists)
	require.NoError(t, err)

	storage.db.Exec(ctx, `TRUNCATE users CASCADE`)
}

func getOpenTasks(ctx context.Context, t *testing.T, storage *Storage) {
	taskOpen1 := &types.Task{
		Url:    randomURL() + fmt.Sprintf("%d", randomId()),
		Status: types.Open,
	}

	err := storage.CreateTask(ctx, taskOpen1)
	require.NoError(t, err)

	taskOpen2 := &types.Task{
		Url:    randomURL() + fmt.Sprintf("%d2", randomId()),
		Status: types.Open,
	}

	err = storage.CreateTask(ctx, taskOpen2)
	require.NoError(t, err)

	taskClosed := &types.Task{
		Url:    randomURL() + fmt.Sprintf("%d3", randomId()),
		Status: types.Closed,
	}

	err = storage.CreateTask(ctx, taskClosed)
	require.NoError(t, err)

	tasks, err := storage.GetOpenTasks(ctx)
	require.NoError(t, err)
	require.Contains(t, tasks, taskOpen1)
	require.Contains(t, tasks, taskOpen2)
	require.NotContains(t, tasks, taskClosed)

	storage.db.Exec(ctx, `TRUNCATE users CASCADE`)
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
