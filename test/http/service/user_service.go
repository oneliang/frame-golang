package service

import (
	"fmt"
)

type UserService struct {
	RoleService *RoleService
}

func NewUserService() *UserService {
	return &UserService{}
}

func (this *UserService) Create() error {
	fmt.Println("UserService.Create")
	return nil
}

func (this *UserService) Read() error {
	fmt.Println("UserService.Read")
	return nil
}

func (this *UserService) Update() error {
	fmt.Println("UserService.Update")
	return nil
}

func (this *UserService) Delete() error {
	fmt.Println("UserService.Delete")
	return nil
}
