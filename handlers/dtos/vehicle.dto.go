package dtos

import "jxb-eprocurement/models"

type VehicleDTO struct {
	ID           uint           `json:"id"`
	Name         string         `json:"name"`
	PoliceNumber string         `json:"police_number"`
	IsAvailable  bool           `json:"is_available"`
	Type         VehicleTypeDTO `json:"type"`
}

func ToVehicleDTO(vehicle models.Vehicle) VehicleDTO {
	return VehicleDTO{
		ID:           vehicle.ID,
		Name:         vehicle.Name,
		PoliceNumber: vehicle.PoliceNumber,
		IsAvailable:  vehicle.IsAvailable,
		Type:         ToVehicleTypeDTO(vehicle.Type),
	}
}
