package service

import (
	"context"
	"backend-service/defn"
)

type UserService struct {
	repo defn.Repository
}

func NewUserService(repo defn.Repository) *UserService {
	return &UserService{repo: repo}
}

func (service *UserService) Create(ctx context.Context, data ...interface{}) (map[string]interface{}, error) {
	return nil, nil
}
