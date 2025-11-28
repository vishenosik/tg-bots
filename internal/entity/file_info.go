package entity

type SenderInfo struct {
	// User info

	UserID    int64
	UserName  string
	FirstName string
	LastName  string

	// Sending info

	ChatID  int64
	PhotoID string
}
