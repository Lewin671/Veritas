package db

import (
	"fmt"
	"log"
	"veritas-server/models"
	"veritas-server/services"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	// Validate encryption key before proceeding
	if err := services.ValidateEncryptionKey(); err != nil {
		log.Fatal("Encryption key validation failed: ", err)
	}
	log.Println("Encryption key validated successfully")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		"localhost",
		"veritas",
		"veritas_password",
		"veritas",
		"5432",
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Database connected successfully")

	// Pre-migration: Add columns without NOT NULL constraint if they don't exist
	DB.Exec("ALTER TABLE model_configs ADD COLUMN IF NOT EXISTS provider text")
	DB.Exec("ALTER TABLE model_configs ADD COLUMN IF NOT EXISTS model_id text")
	
	// Update existing null values with defaults
	DB.Exec("UPDATE model_configs SET provider = 'openai' WHERE provider IS NULL OR provider = ''")
	DB.Exec("UPDATE model_configs SET model_id = 'gpt-4o-mini' WHERE model_id IS NULL OR model_id = ''")

	// Auto Migrate (will add NOT NULL constraints)
	err = DB.AutoMigrate(&models.Conversation{}, &models.Message{}, &models.ModelConfig{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
	log.Println("Database migrated successfully")

	// Run default model config migration
	if err := services.MigrateDefaultModelConfig(DB); err != nil {
		log.Printf("Warning: Failed to migrate default model config: %v", err)
	}
}
