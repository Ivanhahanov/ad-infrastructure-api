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
	Address   string             `bson:"address"`
	SshPubKey string             `bson:"shh_pub_key"`
}

type Flags struct {
	ID          primitive.ObjectID `bson:"_id"`
	Team        string             `bson:"team"`
	Service     string             `bson:"service"`
	AttackFlag  int                `bson:"attack_flag"`
	DefenceFlag int                `bson:"defence_flag"`
}

type JWTTeam struct {
	TeamName string
}