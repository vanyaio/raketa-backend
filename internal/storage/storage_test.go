package storage

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
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
	wg := &sync.WaitGroup{}
	wg.Add(6)
	go func() {
		defer wg.Done()
		err = CreateUser(ctx, t, storage)
		require.NoError(t, err)
	}()
	go func() {
		defer wg.Done()
		err = CreateTask(ctx, t, storage)
		require.NoError(t, err)
	}()
	go func() {
		defer wg.Done()
		err = DeleteTask(ctx, t, storage)
		require.NoError(t, err)
	}()
	go func() {
		defer wg.Done()
		err = AssignUser(ctx, t, storage)
		require.NoError(t, err)
	}()
	go func() {
		defer wg.Done()
		err = CloseTask(ctx, t, storage)
		require.NoError(t, err)
	}()
	go func() {
		defer wg.Done()
		err = GetOpenTasks(ctx, t, storage)
		require.NoError(t, err)
	}()
	wg.Wait()
	
	storage.db.Exec(ctx, `TRUNCATE users CASCADE`)
}

func CreateUser(ctx context.Context, t *testing.T, storage *Storage) error {
	u := &types.User{
		ID: int64(randomId()),
	}

	user, err := storage.CreateUser(ctx, u)
	require.NoError(t, err)
	require.Equal(t, user, u)

	u = &types.User{}
	err = storage.db.QueryRow(ctx, `SELECT * FROM users WHERE id = $1`, user.ID).Scan(&u.ID)
	require.NoError(t, err)
	require.Equal(t, user, u)

	return err
}

func CreateTask(ctx context.Context, t *testing.T, storage *Storage) error {
	task := &types.Task{
		Url:    randomURL(),
		Status: types.Open,
	}

	newTask, err := storage.CreateTask(ctx, task)
	require.NoError(t, err)
	require.Equal(t, newTask, task)

	task = &types.Task{}
	err = storage.db.QueryRow(ctx, `SELECT * FROM tasks WHERE url = $1`, newTask.Url).Scan(&task.Url, &task.UserID, &task.Status)
	require.NoError(t, err)
	require.Equal(t, newTask, task)

	return err
}

func DeleteTask(ctx context.Context, t *testing.T, storage *Storage) error {
	task := &types.Task{
		Url:    randomURL() + fmt.Sprintf("%d", randomId()),
		Status: types.Open,
	}

	newTask, err := storage.CreateTask(ctx, task)
	require.NoError(t, err)
	require.Equal(t, newTask, task)

	err = storage.DeleteTask(ctx, newTask)
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

func AssignUser(ctx context.Context, t *testing.T, storage *Storage) error {
	u := &types.User{
		ID: int64(randomId()),
	}

	user, err := storage.CreateUser(ctx, u)
	require.NoError(t, err)
	require.Equal(t, user, u)

	task := &types.Task{
		Url:    randomURL() + fmt.Sprintf("%d", randomId()),
		Status: types.Open,
	}

	newTask, err := storage.CreateTask(ctx, task)
	require.NoError(t, err)
	require.Equal(t, newTask, task)

	req := &types.AssignUserRequest{
		Url:    newTask.Url,
		UserID: &user.ID,
	}

	newTask, err = storage.AssignUser(ctx, req)
	require.NoError(t, err)
	require.Equal(t, newTask.UserID, &user.ID)

	return err
}

func CloseTask(ctx context.Context, t *testing.T, storage *Storage) error {
	task := &types.Task{
		Url:    randomURL() + fmt.Sprintf("%d", randomId()),
		Status: types.Open,
	}

	newTask, err := storage.CreateTask(ctx, task)
	require.NoError(t, err)
	require.Equal(t, newTask, task)

	req := &types.CloseTaskRequest{
		Url: newTask.Url,
	}

	newTask, err = storage.CloseTask(ctx, req)
	require.NoError(t, err)
	require.NotEqual(t, newTask.Status, task.Status)

	return err
}

func GetOpenTasks(ctx context.Context, t *testing.T, storage *Storage) error {
	task1 := &types.Task{
		Url:    randomURL() + fmt.Sprintf("%d", randomId()),
		Status: types.Open,
	}

	newTask1, err := storage.CreateTask(ctx, task1)
	require.NoError(t, err)
	require.Equal(t, newTask1, task1)

	task2 := &types.Task{
		Url:    randomURL() + fmt.Sprintf("%d2", randomId()),
		Status: types.Open,
	}

	newTask2, err := storage.CreateTask(ctx, task2)
	require.NoError(t, err)
	require.Equal(t, newTask2, task2)

	tasks, err := storage.GetOpenTasks(ctx)
	require.NoError(t, err)
	require.Contains(t, tasks, newTask1)
	require.Contains(t, tasks, newTask2)

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
