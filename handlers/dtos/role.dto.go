package dtos

import "jxb-eprocurement/models"

type RoleDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func ToRoleDTO(role models.Role) RoleDTO {
	return RoleDTO{
		ID:   role.ID,
		Name: role.Name,
	}
}
