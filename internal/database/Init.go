package database

import (
	"OnlineShop/config"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type Customer struct {
	ID               uint   `gorm:"primaryKey"`
	Email            string `gorm:"type:varchar(255);not null;unique"`
	PasswordHash     string `gorm:"type:varchar(255);not null" json:"-"`
	Role             string `gorm:"type:varchar(50);not null;default:'user'"`
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
	Customer   Customer    `gorm:"foreignKey:CustomerID"`
}

type OrderItem struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   uint
	ProductID uint
	Quantity  int
	Price     float64
	Product   Product `gorm:"foreignKey:ProductID"`
}

var DB *gorm.DB

func InitDB(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
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

func CreateInitialAdmin(db *gorm.DB, cfg *config.Config) {
	if cfg.InitialAdminEmail == "" || cfg.InitialAdminPassword == "" {
		log.Println("Initial admin credentials not set, skipping creation.")
		return
	}

	var existingAdmin Customer
	err := db.Where("email = ?", cfg.InitialAdminEmail).First(&existingAdmin).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Creating initial admin user: %s", cfg.InitialAdminEmail)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.InitialAdminPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash initial admin password: %v", err)
		}

		admin := Customer{
			Email:            cfg.InitialAdminEmail,
			PasswordHash:     string(hashedPassword),
			RegistrationDate: time.Now(),
			Role:             "admin",
		}

		if result := db.Create(&admin); result.Error != nil {
			log.Fatalf("Failed to create initial admin: %v", result.Error)
		}
		log.Println("Initial admin created successfully.")
	} else if err != nil {
		log.Fatalf("Failed to query for initial admin: %v", err)
	} else {
		log.Println("Initial admin user already exists.")
	}
}
