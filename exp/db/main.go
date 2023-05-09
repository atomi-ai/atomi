package main

import (
	"fmt"

	"github.com/atomi-ai/atomi/models"
	"github.com/atomi-ai/atomi/repositories"
)

func main() {
	db := models.InitDB()

	userRepo := repositories.NewUserRepository(db)
	user := &models.User{
		Email: "test7@atomi.ai",
		Role:  models.RoleUser,
	}
	_, err := userRepo.Save(user)
	if err != nil {
		fmt.Errorf("Error saving user to database: %w", err)
	}
	fmt.Println("New user ID:", user.ID) // 打印新生成的 ID
}
