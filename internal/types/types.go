package types

type User struct {
	ID int64 `json:"id"`
}

var (
	Open     string = "open"
	Closed          = "closed"
	Declined        = "declined"
)

type Task struct {
	URL    string `json:"url"`
	Status *string `json:"status"`
	UserID *int64 `json:"id"`
}
