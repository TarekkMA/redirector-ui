package data

type Redirect struct {
	From string `gorm:"column:from;primary_key"`
	To   string `gorm:"column:to"`
}
