package database

import (
	"context"
	"errors"
	"github.com/Ivanhahanov/ad-infrastructure-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var collection *mongo.Collection
var flags *mongo.Collection
var ctx = context.TODO()

func InitMongo() {
	adminPass := os.Getenv("ADMIN_PASS")
	if adminPass == "" {
		adminPass = "admin"
	}
	credential := options.Credential{
		Username: "admin",
		Password: adminPass,
	}

	mongoAddr := os.Getenv("MONGODB")
	if mongoAddr == "" {
		mongoAddr = "mongodb://localhost:27017"
	}

	clientOptions := options.Client().ApplyURI(mongoAddr).SetAuth(credential)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("ad").Collection("teams")
	flags = client.Database("ad").Collection("flags")
}

func CreateTeam(team *models.Team) error {
	_, err := collection.InsertOne(ctx, team)
	return err
}
func GetTeams() ([]*models.Team, error) {
	// passing bson.D{{}} matches all documents in the collection
	filter := bson.M{"name": bson.M{"$ne": "admin"}}
	return FilterTeams(filter)
}
func GetUsers() ([]*models.Team, error) {
	// passing bson.D{{}} matches all documents in the collection
	filter := bson.D{{}}
	return FilterTeams(filter)
}

func FilterTeams(filter interface{}) ([]*models.Team, error) {
	// A slice of teams for storing the decoded documents
	var teams []*models.Team

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return teams, err
	}

	for cur.Next(ctx) {
		var t models.Team
		err := cur.Decode(&t)
		if err != nil {
			return teams, err
		}

		teams = append(teams, &t)
	}

	if err := cur.Err(); err != nil {
		return teams, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	if len(teams) == 0 {
		return teams, mongo.ErrNoDocuments
	}

	return teams, nil
}

func DeleteTeam(name string) error {
	filter := bson.D{primitive.E{Key: "name", Value: name}}

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("No teams were deleted")
	}

	return nil
}

func AddAttackFlag(team string, service string) error {
	return AddFlag(team, service, "attack_flag")
}

func AddDefenceFlag(team string, service string) error {
	return AddFlag(team, service, "defence_flag")
}
func AddFlag(team string, service string, field string) error {
	_, err := flags.UpdateOne(ctx, bson.M{
		"team":    team,
		"service": service,
	}, bson.D{
		{"$inc", bson.D{{field, 1}}},
	}, options.Update().SetUpsert(true))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func CountTeams() int64 {
	itemCount, _ := collection.CountDocuments(ctx, bson.M{"name": bson.M{"$ne": "admin"}})
	return itemCount
}
