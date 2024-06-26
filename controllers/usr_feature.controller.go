package controllers

import (
	"jxb-eprocurement/handlers"
	"jxb-eprocurement/service"

	"github.com/gin-gonic/gin"
)

type FeatureController interface {
	GetAllFeatures(c *gin.Context)
	GetFeature(c *gin.Context)
	CreateFeature(c *gin.Context)
	UpdateFeature(c *gin.Context)
	DeleteFeature(c *gin.Context)
}

// FeatureControllerImpl is the implementation of the FeatureController interface.
type FeatureControllerImpl struct {
	service service.FeatureService
}

// NewFeatureController creates a new instance of FeatureControllerImpl.
func FeatureControllerConstructor(service service.FeatureService) FeatureController {
	return &FeatureControllerImpl{service: service}
}

// GetAllFeatures handles the request to get all Features.
func (mc *FeatureControllerImpl) GetAllFeatures(c *gin.Context) {
	response := mc.service.GetAll(c)
	handlers.ResponseFormatterWithLogging(c, response)
}

// GetFeature handles the request to get a Feature by ID.
func (mc *FeatureControllerImpl) GetFeature(c *gin.Context) {
	response := mc.service.GetByID(c)
	if response.Err != nil {
		handlers.ResponseFormatterWithLogging(c, response)
	} else {
		handlers.ResponseFormatterWithLogging(c, response)
	}
}

// CreateFeature handles the request to add a new Feature.
func (mc *FeatureControllerImpl) CreateFeature(c *gin.Context) {
	response := mc.service.AddData(c)
	if response.Err != nil {
		handlers.ResponseFormatterWithLogging(c, response)
	} else {
		handlers.ResponseFormatterWithLogging(c, response)
	}
}

// UpdateFeature handles the request to update a Feature.
func (mc *FeatureControllerImpl) UpdateFeature(c *gin.Context) {
	response := mc.service.UpdateData(c)
	if response.Err != nil {
		handlers.ResponseFormatterWithLogging(c, response)
	} else {
		handlers.ResponseFormatterWithLogging(c, response)
	}
}

// DeleteFeature handles the request to delete a Feature.
func (mc *FeatureControllerImpl) DeleteFeature(c *gin.Context) {
	response := mc.service.DeleteData(c)
	if response.Err != nil {
		handlers.ResponseFormatterWithLogging(c, response)
	} else {
		handlers.ResponseFormatterWithLogging(c, response)
	}
}
