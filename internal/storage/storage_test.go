package storage

import (
	_ "context"
	_ "regexp"
	_ "testing"

	_ "github.com/pashagolub/pgxmock/v2"
	_ "github.com/stretchr/testify/require"
	_ "github.com/vanyaio/raketa-backend/internal/types"
	_ "github.com/vanyaio/raketa-backend/proto"
)

// func Test_CreateUser(t *testing.T) {
// 	t.Parallel()

// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	u := &types.User{
// 		ID: 1234,
// 	}

// 	colums := []string{
// 		"id",
// 	}

// 	row := pgxmock.NewRows(colums).AddRow(
// 		int64(1234),
// 	)

// 	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (id) VALUES ($1) RETURNING *`)).WithArgs(u.ID).WillReturnRows(row)

// 	storage := NewStorage(mock)

// 	user, err := storage.CreateUser(context.Background(), u)
// 	require.NoError(t, err)
// 	require.Equal(t, user, u)
// }

// func Test_CreateTask(t *testing.T) {

// 	t.Parallel()

// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	task := &types.Task{
// 		URL:    "qwerty",
// 		Status: &types.Open,
// 	}

// 	colums := []string{
// 		"url",
// 		"id",
// 		"status",
// 	}

// 	row := pgxmock.NewRows(colums).AddRow(
// 		"qwerty",
// 		nil,
// 		&types.Open,
// 	)

// 	mock.ExpectQuery(
// 		regexp.QuoteMeta(
// 			`INSERT INTO tasks (url, assigned_id, status) VALUES ($1, NULL, $2) RETURNING *`)).WithArgs(task.URL, task.Status).WillReturnRows(row)

// 	storage := NewStorage(mock)

// 	newTask, err := storage.CreateTask(context.Background(), task)
// 	require.NoError(t, err)
// 	require.Equal(t, newTask, task)
// }

// func Test_DeleteTask(t *testing.T) {
// 	t.Parallel()

// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	task := &types.Task{
// 		URL:    "qwerty",
// 		Status: &types.Open,
// 	}

// 	mock.ExpectExec(
// 		regexp.QuoteMeta(
// 			`DELETE FROM tasks WHERE url = $1`)).WithArgs(task.URL).WillReturnResult(pgxmock.NewResult("DELETE", 1))

// 	storage := NewStorage(mock)

// 	err = storage.DeleteTask(context.Background(), task)
// 	require.NoError(t, err)
// }

// func Test_AssignWorker(t *testing.T) {
// 	t.Parallel()

// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()


// 	req := &proto.AssignRequest{
// 		Url:    "qwerty",
// 		UserId: 1234,
// 	}

// 	var id int64 = 1234

// 	task := &types.Task{
// 		URL:    "qwerty",
// 		Status: &types.Open,
// 		UserID: &id,
// 	}

// 	colums := []string{
// 		"url",
// 		"id",
// 		"status",
// 	}

// 	row := pgxmock.NewRows(colums).AddRow(
// 		"qwerty",
// 		&id,
// 		&types.Open,
// 	)

// 	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE tasks
// 			SET assigned_id = COALESCE($1, assigned_id)
// 			WHERE url = $2
// 			RETURNING *`)).WithArgs(req.UserId, req.Url).WillReturnRows(row)

// 	storage := NewStorage(mock)

// 	newTask, err := storage.AssignWorker(context.Background(), req)
// 	require.NoError(t, err)
// 	require.Equal(t, newTask, task)
// }

// func Test_CloseTask(t *testing.T) {
// 	t.Parallel()

// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	req := &proto.CloseRequest{
// 		Url: "qwerty",
// 	}

// 	var id int64 = 1234

// 	task := &types.Task{
// 		URL:    "qwerty",
// 		Status: &types.Closed,
// 		UserID: &id,
// 	}

// 	colums := []string{
// 		"url",
// 		"id",
// 		"status",
// 	}

// 	row := pgxmock.NewRows(colums).AddRow(
// 		"qwerty",
// 		&id,
// 		&types.Closed,
// 	)

// 	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE tasks
// 		SET status = 'closed'
// 		WHERE url = $1
// 		RETURNING *`)).WithArgs(req.Url).WillReturnRows(row)

// 	storage := NewStorage(mock)

// 	newTask, err := storage.CloseTask(context.Background(), req)
// 	require.NoError(t, err)
// 	require.Equal(t, newTask, task)
// }

// func Test_GetOpenTask(t *testing.T) {
// 	t.Parallel()

// 	mock, err := pgxmock.NewPool()
// 	require.NoError(t, err)
// 	defer mock.Close()

// 	var id int64 = 1234

// 	colums1 := []string{
// 		"url",
// 		"id",
// 		"status",
// 	}

// 	row1 := pgxmock.NewRows(colums1).AddRow(
// 		"qwerty",
// 		&id,
// 		&types.Open,
// 	)


// 	colums2 := []string{
// 		"url",
// 		"id",
// 		"status",
// 	}

// 	row2 := pgxmock.NewRows(colums2).AddRow(
// 		"qwerty1",
// 		nil,
// 		&types.Open,
// 	)


// 	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM tasks WHERE status = 'open'`)).WillReturnRows(row1, row2)

// 	storage := NewStorage(mock)

// 	newTask, err := storage.GetOpenTasks(context.Background())
// 	require.NoError(t, err)
// 	require.NotNil(t, newTask)
// }