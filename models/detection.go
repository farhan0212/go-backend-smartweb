package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoundingBox struct {
	X      float64 `bson:"x" json:"x"`
	Y      float64 `bson:"y" json:"y"`
	Width  float64 `bson:"width" json:"width"`
	Height float64 `bson:"height" json:"height"`
}

type Detection struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Label       string             `bson:"label" json:"label"`
	Confidence  float64            `bson:"confidence" json:"confidence"`
	BoundingBox BoundingBox        `bson:"boundingBox" json:"boundingBox"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
}