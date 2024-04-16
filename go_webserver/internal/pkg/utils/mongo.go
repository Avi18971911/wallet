package utils

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ObjectIdToString(insertedId interface{}) (string, error) {
	if oid, ok := insertedId.(primitive.ObjectID); ok {
		res := oid.Hex()
		return res, nil
	}
	err := errors.New("cannot convert object ID to string")
	return "", err
}

func StringToObjectId(stringId string) (primitive.ObjectID, error) {
	objId, err := primitive.ObjectIDFromHex(stringId)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return objId, nil
}
