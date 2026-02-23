package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"url_shorter/internal/entity"
	"url_shorter/internal/repo"
	"url_shorter/internal/service"
	"url_shorter/internal/transport"
	"url_shorter/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	_ "url_shorter/docs"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Используются системные переменные", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	limitStr := os.Getenv("RATE_LIMIT_PER_MIN")
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 10
	}

	var db *gorm.DB
	var dbErr error

	for i := 0; i < 5; i++ {
		db, dbErr = database.NewPostgresDB(dsn)
		if dbErr == nil {
			break
		}
		log.Printf("⏳ Ожидание базы данных... (попытка %d/5)", i+1)
		time.Sleep(5 * time.Second)
	}

	if dbErr != nil {
		log.Fatalf("❌ Не удалось подключиться к БД после 5 попыток: %v", dbErr)
	}

	db.AutoMigrate(&entity.ShortLink{}, &entity.IpLimit{})

	repository := repo.NewRepository(db)

	svc := service.NewService(repository, os.Getenv("BASE_URL"), limit)

	handler := transport.NewHandler(svc)

	app := fiber.New()

	handler.Register(app)
	app.Listen(":8080")
}
