package user

import (
	"github.com/kostinp/edu-platform-backend/internal/user/user"
)

func NewUserService() user.Service {
	repo := user.NewPostgresRepo()
	return user.NewService(repo)
}
