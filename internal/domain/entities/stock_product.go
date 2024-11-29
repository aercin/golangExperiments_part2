package entities

type StockProduct struct {
	Id         int64  `gorm:"column:id;primaryKey"`
	StockId    int64  `gorm:"column:stock_id"`
	ProductId  string `gorm:"column:product_id"`
	Quantity   int    `gorm:"column:quantity"`
	IsModified bool   `gorm:"-"`
}
