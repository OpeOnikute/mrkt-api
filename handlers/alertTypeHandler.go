package handlers

import (
	"mrkt/constants"
	"mrkt/models"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var model models.AlertModel

// AlertTypeHandler ...
type AlertTypeHandler struct{}

// CreateAlertType ...
func (a AlertTypeHandler) CreateAlertType(alertType models.AlertType) (*mongo.InsertOneResult, error) {
	if existingType, _ := model.FindByName(alertType.Name); existingType.Name == alertType.Name {
		var newErr constants.CustomError
		newErr.Msg = constants.ResourceExists("alert type")
		return nil, &newErr
	}
	return model.Create(alertType)
}

// FindByID ...
func (a AlertTypeHandler) FindByID(requestID string) (models.AlertType, error) {
	id, _ := primitive.ObjectIDFromHex(requestID)
	return model.FindByID(id)
}

// DeleteByID ...
func (a AlertTypeHandler) DeleteByID(requestID string) (*mongo.UpdateResult, error) {
	id, _ := primitive.ObjectIDFromHex(requestID)
	return model.DeleteByID(id)
}

// GetMultiple ...
func (a AlertTypeHandler) GetMultiple(query bson.M) ([]models.AlertType, error) {
	return model.FindMany(query)
}

// UpdateByID ...
func (a AlertTypeHandler) UpdateByID(requestID string, alertType models.AlertType) (*mongo.UpdateResult, error) {
	id, _ := primitive.ObjectIDFromHex(requestID)
	return model.UpdateByID(id, alertType)
}
