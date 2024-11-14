package models

type Quotation struct {
	ID    int    `gorm:"column:id;primaryKey;AUTO_INCREMENT"`
	Text  string `gorm:"column:text;type:text"`
	Level string `gorm:"column:level;type:varchar(255)"`
}

func (Quotation) TableName() string {
	return "main"
}
