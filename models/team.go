package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Team struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Name      string             `bson:"name"`
	Hash      string             `bson:"hash"`
}

type ScoreBoard struct {
	ID     primitive.ObjectID `bson:"_id"`
	Name   string             `bson:"name"`
	Attack int                `bson:"attack"`
	SLA    int                `bson:"sla"`
}
