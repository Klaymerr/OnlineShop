package database

import (
	"OnlineShop/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type Customer struct {
	ID               uint   `gorm:"primaryKey"`
	Email            string `gorm:"type:varchar(255);not null;unique"`
	PasswordHash     string `gorm:"type:varchar(255);not null" json:"-"`
	RegistrationDate time.Time
	Orders           []Order
}

type Product struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"type:varchar(255);not null"`
	Price float64
}

type Order struct {
	ID         uint `gorm:"primaryKey"`
	CustomerID uint
	OrderDate  time.Time
	Status     string      `gorm:"type:varchar(50);not null"`
	Items      []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID              uint `gorm:"primaryKey"`
	OrderID         uint
	ProductID       uint
	Quantity        int
	PriceAtPurchase float64
	Product         Product `gorm:"foreignKey:ProductID"`
}

var DB *gorm.DB

func InitDB(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	err = DB.AutoMigrate(&Product{}, &Customer{}, &Order{}, &OrderItem{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
}
