package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Flags struct {
	ID          primitive.ObjectID `bson:"_id"`
	Team        string             `bson:"team"`
	Service     string             `bson:"service"`
	AttackFlag  int                `bson:"attack_flag"`
	DefenceFlag int                `bson:"defence_flag"`
}
