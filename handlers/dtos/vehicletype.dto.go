package dtos

import "jxb-eprocurement/models"

type VehicleTypeDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func ToVehicleTypeDTO(vehicleType models.VehicleType) VehicleTypeDTO {
	return VehicleTypeDTO{
		ID:   vehicleType.ID,
		Name: vehicleType.Name,
	}
}
