package initializers

import "github.com/AdrianTworek/go-tasks-manager/models"

func SyncDb() {
	DB.AutoMigrate(&models.User{})
}
