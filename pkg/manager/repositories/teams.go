package repositories

import (
	"errors"
	"github.com/Ivanhahanov/ad-infrastructure-api/internal/entity"
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TeamRepository struct {db *db.Db}

func NewTeamRepository(db *db.Db) *TeamRepository {
	return &TeamRepository{db: db}
}

func (t *TeamRepository) CreateTeam(team *entity.Team) error {
	_, err := t.db.Teams.InsertOne(t.db.Ctx, team)

	return err
}

func (t *TeamRepository) GetTeams() ([]*entity.Team, error) {return t.FilterTeams(bson.M{"name": bson.M{"$ne": "admin"}})}
func (t *TeamRepository) GetUsers() ([]*entity.Team, error) {return t.FilterTeams(bson.D{{}})}

func (t *TeamRepository) FilterTeams(filter interface{}) ([]*entity.Team, error) {
	var teams []*entity.Team

	cur, err := t.db.Teams.Find(t.db.Ctx, filter)
	if err != nil {return teams, err}

	for cur.Next(t.db.Ctx) {
		var t entity.Team

		err := cur.Decode(&t)
		if err != nil {return teams, err}

		teams = append(teams, &t)
	}

	if err := cur.Err(); err != nil {return teams, err}

	cur.Close(t.db.Ctx)

	if len(teams) == 0 {return teams, mongo.ErrNoDocuments}

	return teams, nil
}

func (t *TeamRepository) DeleteTeam(name string) error {
	res, err := t.db.Teams.DeleteOne(t.db.Ctx, bson.D{primitive.E{Key: "name", Value: name}})
	if err != nil {return err}

	if res.DeletedCount == 0 {return errors.New("No teams were deleted")}

	return nil
}

