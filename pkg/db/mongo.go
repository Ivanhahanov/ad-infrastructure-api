package db

import (
	"context"
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Db struct {
	Teams *mongo.Collection
	Flags *mongo.Collection

	Ctx context.Context
}

func New() *Db {
	addr := utils.GetEnv("MONGODB", "mongodb://localhost:27017")

	clientOptions := options.Client().ApplyURI(addr).SetAuth(options.Credential{
		Username: utils.GetEnv("MONGO_USER", "admin"),
		Password: utils.GetEnv("ADMIN_PASS", "admin"),
	})

	ctx := context.TODO()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {log.Fatal(err)}

	if err = client.Ping(ctx, nil); err != nil {log.Fatal(err)}

	db := client.Database("ad")

	return &Db{
		Teams: db.Collection("teams"),
		Flags: db.Collection("flags"),

		Ctx: ctx,
	}
}

