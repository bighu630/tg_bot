package models

type Quotation struct {
	ID    int    `gorm:"column:id;primaryKey;AUTO_INCREMENT"`
	Data  string `gorm:"column:data;type:text;not null"`
	Level string `gorm:"column:level;type:varchar(255)"`
}

func (Quotation) TableName() string {
	return "quotation"
}
