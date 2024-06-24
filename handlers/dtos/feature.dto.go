package dtos

import "jxb-eprocurement/models"

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
	ModuleID uint   `json:"module_id" form:"module_id" validate:"required,is-uint"`
}

// // ToUSRFeatureDTO converts a USR_Feature model to a USRFeatureDTO in detail format.
// // This function recursively converts child modules as well.
// // Use this function to where child model is needed.
// func ToUSRFeatureDTO(module models.USR_Feature) USRFeatureDTO {
// 	// Convert child modules
// 	children := make([]USRFeatureDTO, len(module.Child))
// 	for i, child := range module.Child {
// 		children[i] = ToUSRFeatureDTO(child)
// 	}

// 	// Return the DTO with converted fields
// 	return USRFeatureDTO{
// 		ID:       module.ID,
// 		Name:     module.Name,
// 		ModuleID: module.ModuleID,
// 		Children: children,
// 	}
// }

// // ToUSRFeatureModel converts a USRFeatureDTO to a USR_Feature model in detail format.
// // This function recursively converts child DTOs as well.
// // Use this function to where child model is needed.
// func ToUSRFeatureModel(dto USRFeatureDTO) models.USR_Feature {
// 	// Convert child DTOs
// 	children := make([]models.USR_Feature, len(dto.Children))
// 	for i, child := range dto.Children {
// 		children[i] = ToUSRFeatureModel(child)
// 	}

// 	// Return the model with converted fields
// 	return models.USR_Feature{
// 		ID:       dto.ID,
// 		Name:     dto.Name,
// 		ModuleID: dto.ModuleID,
// 		Child:    children,
// 	}
// }

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

// // ToUSRFeatureDTO converts a USR_Feature model to a USRFeatureDTO in minimal format.
// // Use this function to where child model not needed.
// func ToUSRFeatureMinimalDTOs(modules []models.USR_Feature) []USRFeatureMinimalDTO {
// 	var moduleDTOs []USRFeatureMinimalDTO

// 	for _, module := range modules {
// 		fmt.Println(module)

// 		moduleDTOs = append(moduleDTOs, ToUSRFeatureMinimalDTO(module))
// 	}

// 	// Return the DTO with converted fields
// 	return moduleDTOs
// }

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
