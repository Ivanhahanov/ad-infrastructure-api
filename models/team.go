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

//type ScoreBoard struct {
//	ID     primitive.ObjectID `bson:"_id"`
//	Name   string             `bson:"name"`
//	Attack int                `bson:"attack"`
//	SLA    int                `bson:"sla"`
//}

type JWTTeam struct {
	TeamName string
}
type Scoreboard struct {
	Teams []ScoreboardTeam
}

type ScoreboardTeam struct {
	TeamName string
	Services []ScoreboardService
}

type ScoreboardService struct {
	ServiceName string
	Routes      map[string]float64
}
