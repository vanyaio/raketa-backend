package storage

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/vanyaio/raketa-backend/config"
	"github.com/vanyaio/raketa-backend/internal/types"
	"github.com/vanyaio/raketa-backend/pkg/db"
)

type RaketaTestSuite struct {
	suite.Suite
	db      *pgxpool.Pool
	storage *Storage
	ctx     context.Context
}

func TestRaketaTestSuit(t *testing.T) {
	suite.Run(t, &RaketaTestSuite{})
}

func (r *RaketaTestSuite) SetupSuite() {
	config := config.GetConfig()

	ctx := context.Background()

	pool, err := db.NewPool(ctx, config)
	r.NoError(err)

	storage := NewStorage(pool)

	r.db = pool
	r.storage = storage
	r.ctx = ctx
}

func (r *RaketaTestSuite) TearDownSuite() {
	r.db.Close()
}

func (r *RaketaTestSuite) TearDownSubTest() {
	r.db.Exec(r.ctx, `TRUNCATE users CASCADE`)
}

func (r *RaketaTestSuite) TearDownTest() {
	r.db.Exec(r.ctx, `TRUNCATE users CASCADE`)
}

func (r *RaketaTestSuite) Test_CreateUser() {
	u := &types.User{
		ID:       int64(randomId()),
		Username: randomURL(),
	}

	err := r.storage.CreateUser(r.ctx, u)
	r.NoError(err)

	var exists bool
	query := `SELECT EXISTS (SELECT * FROM users WHERE id = $1 AND telegram_username = $2)`
	err = r.db.QueryRow(r.ctx, query, u.ID, u.Username).Scan(&exists)
	if err != nil {
		r.Error(err)
	}
	r.True(exists)
	r.NoError(err)
}

func (r *RaketaTestSuite) Test_CreateTask() {
	task := &types.Task{
		Url:    randomURL(),
		Status: types.Open,
	}

	err := r.storage.CreateTask(r.ctx, task)
	r.NoError(err)

	var exists bool
	query := `SELECT EXISTS (SELECT * FROM tasks WHERE url = $1)`
	err = r.storage.db.QueryRow(r.ctx, query, task.Url).Scan(&exists)
	if err != nil {
		r.Error(err)
	}
	r.True(exists)
	r.NoError(err)
}

func (r *RaketaTestSuite) Test_DeleteTask() {
	task := &types.Task{
		Url:    randomURL(),
		Status: types.Open,
	}

	err := r.storage.CreateTask(r.ctx, task)
	r.NoError(err)

	err = r.storage.DeleteTask(r.ctx, task)
	r.NoError(err)

	var notExists bool
	query := `SELECT NOT EXISTS (SELECT * FROM tasks WHERE url = $1)`
	err = r.storage.db.QueryRow(r.ctx, query, task.Url).Scan(&notExists)
	if err != nil {
		r.Error(err)
	}
	r.True(notExists)
	r.NoError(err)
}

func (r *RaketaTestSuite) Test_AssignUser() {
	r.Run("user exists", func() {
		u := &types.User{
			ID:       int64(randomId()),
			Username: randomURL(),
		}

		err := r.storage.CreateUser(r.ctx, u)
		r.NoError(err)

		task := &types.Task{
			Url:    randomURL(),
			Status: types.Open,
		}

		err = r.storage.CreateTask(r.ctx, task)
		r.NoError(err)

		req := &types.AssignUserRequest{
			Url:      task.Url,
			Username: u.Username,
		}

		err = r.storage.AssignUser(r.ctx, req)
		r.NoError(err)

		var exists bool
		query := `SELECT EXISTS (SELECT * FROM tasks WHERE url = $1 AND assigned_id IS NOT NULL)`
		err = r.storage.db.QueryRow(r.ctx, query, task.Url).Scan(&exists)
		if err != nil {
			r.Error(err)
		}
		r.True(exists)
		r.NoError(err)
	})

	r.Run("user doesn't exist", func() {
		u := &types.User{
			ID:       int64(randomId()),
			Username: randomURL(),
		}

		task := &types.Task{
			Url:    randomURL(),
			Status: types.Open,
		}

		err := r.storage.CreateTask(r.ctx, task)
		r.NoError(err)

		req := &types.AssignUserRequest{
			Url:      task.Url,
			Username: u.Username,
		}

		err = r.storage.AssignUser(r.ctx, req)
		r.EqualError(err, pgx.ErrNoRows.Error())
	})
}

func (r *RaketaTestSuite) Test_CloseTask() {
	task := &types.Task{
		Url:    randomURL(),
		Status: types.Open,
	}

	err := r.storage.CreateTask(r.ctx, task)
	r.NoError(err)

	req := &types.CloseTaskRequest{
		Url: task.Url,
	}

	err = r.storage.CloseTask(r.ctx, req)
	r.NoError(err)

	var exists bool
	query := `SELECT EXISTS (SELECT * FROM tasks WHERE url = $1 AND status = 'closed')`
	err = r.storage.db.QueryRow(r.ctx, query, req.Url).Scan(&exists)
	if err != nil {
		r.Error(err)
	}
	r.True(exists)
	r.NoError(err)
}

func (r *RaketaTestSuite) Test_GetUnassignTasks() {
	u := &types.User{
		ID:       int64(randomId()),
		Username: randomURL(),
	}

	err := r.storage.CreateUser(r.ctx, u)
	r.NoError(err)


	taskOpen1 := &types.Task{
		Url:    randomURL() + fmt.Sprintf("%d", randomId()),
		Status: types.Open,
	}

	err = r.storage.CreateTask(r.ctx, taskOpen1)
	r.NoError(err)

	req := &types.AssignUserRequest{
		Url:      taskOpen1.Url,
		Username: u.Username,
	}

	err = r.storage.AssignUser(r.ctx, req)
	r.NoError(err)

	taskOpen2 := &types.Task{
		Url:    randomURL() + fmt.Sprintf("%d2", randomId()),
		Status: types.Open,
	}

	err = r.storage.CreateTask(r.ctx, taskOpen2)
	r.NoError(err)

	taskClosed := &types.Task{
		Url:    randomURL() + fmt.Sprintf("%d3", randomId()),
		Status: types.Closed,
	}

	err = r.storage.CreateTask(r.ctx, taskClosed)
	r.NoError(err)

	tasks, err := r.storage.GetUnassignTasks(r.ctx)
	r.NoError(err)
	r.NotContains(tasks, taskOpen1)
	r.Contains(tasks, taskOpen2)
	r.NotContains(tasks, taskClosed)
}

func (r *RaketaTestSuite) Test_SetTaskPrice() {
	task := &types.Task{
		Url:    randomURL(),
		Status: types.Open,
	}

	err := r.storage.CreateTask(r.ctx, task)
	r.NoError(err)

	req := &types.SetTaskPriceRequest{
		Url:   task.Url,
		Price: uint64(randomId()),
	}

	err = r.storage.SetTaskPrice(r.ctx, req)
	r.NoError(err)

	var exists bool
	query := `SELECT EXISTS (SELECT * FROM tasks WHERE url = $1 AND price = $2)`
	err = r.storage.db.QueryRow(r.ctx, query, req.Url, req.Price).Scan(&exists)
	if err != nil {
		r.Error(err)
	}
	r.True(exists)
	r.NoError(err)
}

func (r *RaketaTestSuite) Test_GetUserTasks() {
	u := &types.User{
		ID:       int64(randomId()),
		Username: randomURL(),
	}

	err := r.storage.CreateUser(r.ctx, u)
	r.NoError(err)

	taskClosed := &types.Task{
		Url:    randomURL() + fmt.Sprintf("1"),
		Status: types.Open,
	}

	taskOpen := &types.Task{
		Url:    randomURL() + fmt.Sprintf("2"),
		Status: types.Open,
	}

	err = r.storage.CreateTask(r.ctx, taskClosed)
	r.NoError(err)

	err = r.storage.CreateTask(r.ctx, taskOpen)
	r.NoError(err)

	reqClosed := &types.AssignUserRequest{
		Url:      taskClosed.Url,
		Username: u.Username,
	}

	reqOpen := &types.AssignUserRequest{
		Url:      taskOpen.Url,
		Username: u.Username,
	}

	err = r.storage.AssignUser(r.ctx, reqClosed)
	r.NoError(err)

	err = r.storage.AssignUser(r.ctx, reqOpen)
	r.NoError(err)

	err = r.storage.CloseTask(r.ctx, &types.CloseTaskRequest{
		Url: reqClosed.Url,
	})

	tasksCount, err := r.storage.GetUserStats(r.ctx, u)
	r.NoError(err)
	r.Equal(tasksCount, int64(1))
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
	return rand.Intn(1000000000) + 1000000000
}
