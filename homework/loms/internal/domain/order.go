package domain

type Order struct {
	ID     int64
	Status string
	UserID int64
	Items  []Item
}
