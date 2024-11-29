package persistence

import (
	"fmt"
	"log"

	"go-poc/configs/abstractions"
	"go-poc/internal/domain/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDbConnection(configs abstractions.Config) *gorm.DB {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		configs.GetValue("Postgres:Host"),
		configs.GetValue("Postgres:Port"),
		configs.GetValue("Postgres:UserName"),
		configs.GetValue("Postgres:Password"),
		configs.GetValue("Postgres:DatabaseName"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err)
	}

	// Tabloyu oluştur (eğer yoksa)
	if err := db.AutoMigrate(&entities.Stock{}, &entities.StockProduct{}, &entities.OutboxMessage{}, &entities.InboxMessage{}); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	return db
}
