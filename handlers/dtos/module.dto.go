package dtos

import (
	"fmt"
	"jxb-eprocurement/models"
)

// USRModuleDTO represents a Data Transfer Object for the USR_Module model in detail format.
// It includes only the fields necessary for data transfer and serialization.
type USRModuleDTO struct {
	ID       uint           `json:"id"`        // Unique identifier of the module
	Name     string         `json:"name"`      // Name of the module
	ParentID *uint          `json:"parent_id"` // ID of the parent module, if any
	Children []USRModuleDTO `json:"children"`  // Child modules
}

// USRModuleDTO represents a Data Transfer Object for the USR_Module model in minimal format.
// It includes only the fields necessary for data transfer and serialization.
type USRModuleMinimalDTO struct {
	ID       uint   `json:"id"`        // Unique identifier of the module
	Name     string `json:"name"`      // Name of the module
	ParentID *uint  `json:"parent_id"` // ID of the parent module, if any
}

// ToUSRModuleDTO converts a USR_Module model to a USRModuleDTO in detail format.
// This function recursively converts child modules as well.
// Use this function to where child model is needed.
func ToUSRModuleDTO(module models.USR_Module) USRModuleDTO {
	// Convert child modules
	children := make([]USRModuleDTO, len(module.Child))
	for i, child := range module.Child {
		children[i] = ToUSRModuleDTO(child)
	}

	// Return the DTO with converted fields
	return USRModuleDTO{
		ID:       module.ID,
		Name:     module.Name,
		ParentID: module.ParentID,
		Children: children,
	}
}

// ToUSRModuleModel converts a USRModuleDTO to a USR_Module model in detail format.
// This function recursively converts child DTOs as well.
// Use this function to where child model is needed.
func ToUSRModuleModel(dto USRModuleDTO) models.USR_Module {
	// Convert child DTOs
	children := make([]models.USR_Module, len(dto.Children))
	for i, child := range dto.Children {
		children[i] = ToUSRModuleModel(child)
	}

	// Return the model with converted fields
	return models.USR_Module{
		ID:       dto.ID,
		Name:     dto.Name,
		ParentID: dto.ParentID,
		Child:    children,
	}
}

// ToUSRModuleDTO converts a USR_Module model to a USRModuleDTO in minimal format.
// Use this function to where child model not needed.
func ToUSRModuleMinimalDTO(module models.USR_Module) USRModuleMinimalDTO {
	// Return the DTO with converted fields
	return USRModuleMinimalDTO{
		ID:       module.ID,
		Name:     module.Name,
		ParentID: module.ParentID,
	}
}

// ToUSRModuleDTO converts a USR_Module model to a USRModuleDTO in minimal format.
// Use this function to where child model not needed.
func ToUSRModuleMinimalDTOs(modules []models.USR_Module) []USRModuleMinimalDTO {
	var moduleDTOs []USRModuleMinimalDTO

	for _, module := range modules {
		fmt.Println(module)

		moduleDTOs = append(moduleDTOs, ToUSRModuleMinimalDTO(module))
	}

	// Return the DTO with converted fields
	return moduleDTOs
}

// ToUSRModuleModel converts a USRModuleDTO to a USR_Module model in minimal format.
// Use this function to where child model not needed.
func ToUSRModuleMinimalModel(dto USRModuleMinimalDTO) models.USR_Module {
	// Return the model with converted fields
	return models.USR_Module{
		ID:       dto.ID,
		Name:     dto.Name,
		ParentID: dto.ParentID,
	}
}
