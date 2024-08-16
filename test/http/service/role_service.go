package service

import (
	"fmt"
)

type RoleService struct {
}

func NewRoleService() *RoleService {
	return &RoleService{}
}

func (this *RoleService) Create() error {
	fmt.Println("RoleService.Create")
	return nil
}

func (this *RoleService) Read() error {
	fmt.Println("RoleService.Read")
	return nil
}

func (this *RoleService) Update() error {
	fmt.Println("RoleService.Update")
	return nil
}

func (this *RoleService) Delete() error {
	fmt.Println("RoleService.Delete")
	return nil
}
