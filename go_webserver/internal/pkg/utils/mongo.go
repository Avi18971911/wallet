package utils

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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

func GetCurrentTimestamp() primitive.Timestamp {
	return TimeToTimestamp(time.Now())
}
func TimeToTimestamp(tm time.Time) primitive.Timestamp {
	return primitive.Timestamp{T: uint32(tm.Unix())}
}

func TimestampToTime(ts primitive.Timestamp) time.Time { return time.Unix(int64(ts.T), 0) }
