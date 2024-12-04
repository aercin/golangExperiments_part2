package entities

type Stock struct {
	Id            int64           `gorm:"column:id;primaryKey"`
	StockProducts []*StockProduct `gorm:"foreignKey:StockId;references:Id"`
}
