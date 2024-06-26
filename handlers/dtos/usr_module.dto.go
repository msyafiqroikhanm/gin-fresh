package dtos

import (
	"jxb-eprocurement/models"
)

// USRModuleDTO represents a Data Transfer Object for the USR_Module model in detail format.
// It includes only the fields necessary for data transfer and serialization.
type (
	USRModuleDTO struct {
		ID       uint           `json:"id"`        // Unique identifier of the module
		Name     string         `json:"name"`      // Name of the module
		ParentID *uint          `json:"parent_id"` // ID of the parent module, if any
		Children []USRModuleDTO `json:"children"`  // Child modules
	}

	// USRModuleDTO represents a Data Transfer Object for the USR_Module model in minimal format.
	// It includes only the fields necessary for data transfer and serialization.
	USRModuleMinimalDTO struct {
		ID       uint   `json:"id" form:"id"`                         // Unique identifier of the module
		Name     string `json:"name" form:"name" validate:"required"` // Name of the module
		ParentID *uint  `json:"parent_id" form:"parent_id"`           // ID of the parent module, if any
	}

	USRModuleFeaturesDTO struct {
		ID       uint                   `json:"id" form:"id"`                         // Unique identifier of the module
		Name     string                 `json:"name" form:"name" validate:"required"` // Name of the module
		ParentID *uint                  `json:"parent_id" form:"parent_id"`           // ID of the parent module, if any
		Children []USRModuleMinimalDTO  `json:"children"`                             // Child modules
		Features []USRFeatureMinimalDTO `json:"features"`                             // List of feature
	}
)

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

func ToUSRModuleWithFeaturesDTO(module models.USR_Module) USRModuleFeaturesDTO {
	// Convert child modules
	children := ToUSRModuleMinimalDTOs(module.Child)
	if len(children) == 0 {
		children = []USRModuleMinimalDTO{}
	}

	// Convert Feature
	features := ToUSRFeatureMinimalDTOs(module.Features)
	if len(features) == 0 {
		features = []USRFeatureMinimalDTO{}
	}

	// Return the DTO with converted fields
	return USRModuleFeaturesDTO{
		ID:       module.ID,
		Name:     module.Name,
		ParentID: module.ParentID,
		Children: children,
		Features: features,
	}
}

func MinimalUSRModuleDTOToInterfaceSlice(slice []USRModuleMinimalDTO) []interface{} {
	interfaceSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}
