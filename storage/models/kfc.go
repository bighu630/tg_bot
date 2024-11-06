package models

type GroupData struct {
	KFC string `gorm:"column:kfc;type:text;comment:KFC"`
}

func (GroupData) TableName() string {
	return "chat"
}
