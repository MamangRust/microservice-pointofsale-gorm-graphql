package repository

import "gorm.io/gorm"

type Repositories struct {
	User         UserRepository
	RefreshToken RefreshTokenRepository
	UserRole     UserRoleRepository
	Role         RoleRepository
	ResetToken   ResetTokenRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:         NewUserRepository(db),
		RefreshToken: NewRefreshTokenRepository(db),
		UserRole:     NewUserRoleRepository(db),
		Role:         NewRoleRepository(db),
		ResetToken:   NewResetTokenRepository(db),
	}
}
