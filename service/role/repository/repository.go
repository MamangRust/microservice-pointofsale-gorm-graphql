package repository

import "gorm.io/gorm"

type Repositories struct {
	RoleCommand RoleCommandRepository
	RoleQuery   RoleQueryRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		RoleCommand: NewRoleCommandRepository(db),
		RoleQuery:   NewRoleQueryRepository(db),
	}
}
