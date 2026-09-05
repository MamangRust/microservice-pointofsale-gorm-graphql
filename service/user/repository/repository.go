package repository

import "gorm.io/gorm"

type Repositories struct {
	UserCommand UserCommandRepository
	UserQuery   UserQueryRepository
	Role        RoleQueryRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		UserCommand: NewUserCommandRepository(db),
		UserQuery:   NewUserQueryRepository(db),
		Role:        NewRoleRepository(db),
	}
}
