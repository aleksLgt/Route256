package domain

type ListItem struct {
	SKU   int64
	Count uint16
	Name  string
	Price uint32
}

type Item struct {
	SKU   int64
	Count uint16
}
