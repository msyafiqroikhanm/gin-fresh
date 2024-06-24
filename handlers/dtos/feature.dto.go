package dtos

import (
	"jxb-eprocurement/models"
)

// USRFeatureDTO represents a Data Transfer Object for the USR_Feature model in detail format.
// It includes only the fields necessary for data transfer and serialization.
type USRFeatureDTO struct {
	ID       uint            `json:"id"`        // Unique identifier of the module
	Name     string          `json:"name"`      // Name of the module
	ModuleID uint            `json:"parent_id"` // ID of the parent module, if any
	Children []USRFeatureDTO `json:"children"`  // Child modules
}

// USRFeatureDTO represents a Data Transfer Object for the USR_Feature model in minimal format.
// It includes only the fields necessary for data transfer and serialization.
type USRFeatureMinimalDTO struct {
	ID       uint   `json:"id" form:"id"`                         // Unique identifier of the module
	Name     string `json:"name" form:"name" validate:"required"` // Name of the module
	ModuleID uint   `json:"module_id" form:"module_id" validate:"required"`
}

type USRFeatureWithModuleDTO struct {
	ID     uint         `json:"id" form:"id"`                         // Unique identifier of the module
	Name   string       `json:"name" form:"name" validate:"required"` // Name of the module
	Module USRModuleDTO `json:"module"`
}

// ToUSRFeatureDTO converts a USR_Feature model to a USRFeatureDTO in minimal format.
// Use this function to where child model not needed.
func ToUSRFeatureMinimalDTO(module models.USR_Feature) USRFeatureMinimalDTO {
	// Return the DTO with converted fields
	return USRFeatureMinimalDTO{
		ID:       module.ID,
		Name:     module.Name,
		ModuleID: module.ModuleID,
	}
}

// ToUSRFeatureDTO converts a USR_Feature model to a USRFeatureDTO in minimal format.
// Use this function to where child model not needed.
func ToUSRFeatureMinimalDTOs(modules []models.USR_Feature) []USRFeatureMinimalDTO {
	var moduleDTOs []USRFeatureMinimalDTO

	for _, module := range modules {
		// fmt.Println(module)

		moduleDTOs = append(moduleDTOs, ToUSRFeatureMinimalDTO(module))
	}

	// Return the DTO with converted fields
	return moduleDTOs
}

func ToUSRFeatureMinimalWithModuleDTO(feature models.USR_Feature) USRFeatureWithModuleDTO {
	// Return the DTO with converted fields
	return USRFeatureWithModuleDTO{
		ID:     feature.ID,
		Name:   feature.Name,
		Module: ToUSRModuleDTO(feature.Module),
	}
}

func ToUSRFeatureMinimalWithModuleDTOs(features []models.USR_Feature) []USRFeatureWithModuleDTO {
	var featureDTOs []USRFeatureWithModuleDTO

	for _, feature := range features {
		featureDTOs = append(featureDTOs, ToUSRFeatureMinimalWithModuleDTO(feature))
	}

	// Memeriksa apakah slice featureDTOs kosong
	if len(featureDTOs) == 0 {
		// Jika kosong, mengembalikan sebuah slice kosong dengan tipe yang sesuai
		return []USRFeatureWithModuleDTO{}
	}

	// Jika tidak kosong, kembalikan featureDTOs yang telah diisi dengan data
	return featureDTOs
}

// ToUSRFeatureModel converts a USRFeatureDTO to a USR_Feature model in minimal format.
// Use this function to where child model not needed.
func ToUSRFeatureMinimalModel(dto USRFeatureMinimalDTO) models.USR_Feature {
	// Return the model with converted fields
	return models.USR_Feature{
		ID:       dto.ID,
		Name:     dto.Name,
		ModuleID: dto.ModuleID,
	}
}
