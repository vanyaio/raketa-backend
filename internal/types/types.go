package types

type User struct {
	ID int64 `json:"id"`
}

type Status string

var (
	Open     Status = "open"
	Closed   Status = "closed"
	Declined Status = "declined"
)

type Task struct {
	Url    string `json:"url"`
	Status Status `json:"status"`
	UserID *int64 `json:"id"`
}

type AssignUserRequest struct {
	Url    string `json:"url"`
	UserID *int64 `json:"id"`
}

type CloseTaskRequest struct {
	Url string `json:"url"`
}
