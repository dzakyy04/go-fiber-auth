package database

import (
	"fmt"
	"go-fiber-auth/models/entity"
)

func MigrateDatabase() {
	err := DB.AutoMigrate(&entity.User{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	fmt.Println("Successfully migrated the database.")
}
