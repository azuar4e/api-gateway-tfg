package initializers

import (
	"github.com/azuar4e/api-gateway-tfg/internal/models"
)

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}
