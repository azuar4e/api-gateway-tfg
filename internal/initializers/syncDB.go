package initializers

import (
	"github.com/azuar4e/microservices-tfg/internal/models"
)

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}
