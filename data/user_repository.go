package data

import (
	"context"
)

type UserRepository struct {
	db Database
}

func NewUserRepository(db Database) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) Create(ctx context.Context, data ...interface{}) (map[string]interface{}, error) {
	return nil, nil
}
