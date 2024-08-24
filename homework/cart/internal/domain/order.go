package domain

type Order struct {
	Status string
	UserID int64
	Items  []Item
}
