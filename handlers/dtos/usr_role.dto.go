package dtos

import (
	"jxb-eprocurement/models"
)

type (
	// USRRoleDTO represents a Data Transfer Object for the USR_Role model in detail format.
	// It includes only the fields necessary for data transfer and serialization.
	// TODO: Adding list of feature that the role have
	USRRoleDTO struct {
		ID               uint   `json:"id" form:"id"`
		Name             string `json:"name" form:"name"`
		IsAdministrative bool   `json:"is_administrative" form:"is_administrative"`
		Features         []uint `json:"features" form:"features"`
	}

	// USRRoleDTO represents a Data Transfer Object for the USR_Role model in minimal format.
	// It includes only the fields necessary for data transfer and serialization.
	USRRoleMinimalDTO struct {
		ID               uint   `json:"id" form:"id"`
		Name             string `json:"name" form:"name" validate:"required"`
		IsAdministrative bool   `json:"is_administrative"`
	}

	// DTO that serialization input from user for method POST and PUT
	InputUSRRoleDTO struct {
		Name             string `json:"name" form:"name" validate:"required"`
		IsAdministrative bool   `json:"is_administrative" form:"is_administrative" validate:"boolean"`
		Features         []uint `json:"features" form:"features" validate:"required"`
	}
)

// ToUSRRoleDTO converts a USR_Role model to a USRRoleDTO in detail format.
// Use this function to where the feature that role have is needed.
func ToUSRRoleDTO(role models.USR_Role) USRRoleDTO {
	var featureIDs []uint
	for _, feature := range role.Features {
		featureIDs = append(featureIDs, feature.ID)
	}

	return USRRoleDTO{
		ID:               role.ID,
		Name:             role.Name,
		IsAdministrative: role.IsAdministrative,
		Features:         featureIDs,
	}
}

// ToUSRRoleModel converts a USRRoleDTO to a USR_Role model in detail format.
// Use this function to where the feature that role have is needed.
func ToUSRRoleModel(dto USRRoleDTO) models.USR_Role {
	// Return the model with converted fields
	return models.USR_Role{
		ID:               dto.ID,
		Name:             dto.Name,
		IsAdministrative: dto.IsAdministrative,
	}
}

// InputToUSRRoleModel converts a serialization from InputUSRRoleDTO to a USR_Role model in detail format.
// Use this function to where the feature that role have is needed.
func InputToUSRRoleModel(dto InputUSRRoleDTO) models.USR_Role {
	// Return the model with converted fields
	return models.USR_Role{
		Name:             dto.Name,
		IsAdministrative: dto.IsAdministrative,
	}
}

// ToUSRRoleDTO converts a USR_Role model to a USRRoleDTO in minimal format.
// Use this function to where feature model not needed.
func ToUSRRoleMinimalDTO(role models.USR_Role) USRRoleMinimalDTO {
	// Return the DTO with converted fields
	return USRRoleMinimalDTO{
		ID:               role.ID,
		Name:             role.Name,
		IsAdministrative: role.IsAdministrative,
	}
}

// ToUSRRoleDTO converts a USR_Role model to a USRRoleDTO in minimal format.
// Use this function to where feature model not needed.
func ToUSRRoleMinimalDTOs(roles []models.USR_Role) []USRRoleMinimalDTO {
	var roleDTOs []USRRoleMinimalDTO

	for _, role := range roles {
		roleDTOs = append(roleDTOs, ToUSRRoleMinimalDTO(role))
	}

	if len(roleDTOs) == 0 {
		roleDTOs = []USRRoleMinimalDTO{}
	}

	// Return the DTO with converted fields
	return roleDTOs
}

// ToUSRRoleModel converts a USRRoleDTO to a USR_Role model in minimal format.
// Use this function to where feature model not needed.
func ToUSRRoleMinimalModel(dto USRRoleMinimalDTO) models.USR_Role {
	// Return the model with converted fields
	return models.USR_Role{
		ID:               dto.ID,
		Name:             dto.Name,
		IsAdministrative: dto.IsAdministrative,
	}
}

func MinimalRoleDTOToInterfaceSlice(slice []USRRoleMinimalDTO) []interface{} {
	interfaceSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}
