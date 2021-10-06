package repositories

import (
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type FlagRepository struct {db *db.Db}

func NewFlagRepository(db *db.Db) *FlagRepository {
	return &FlagRepository{db: db}
}

func (r *FlagRepository) AddAttackFlag (team string, service string) error {return r.addFlag(team, service, "attack_flag")}
func (r *FlagRepository) AddDefenceFlag(team string, service string) error {return r.addFlag(team, service, "defence_flag")}

func (r *FlagRepository) addFlag(team string, service string, field string) error {
	_, err := r.db.Flags.UpdateOne(r.db.Ctx, bson.M{
		"team":    team,
		"service": service,
	}, bson.D{
		{"$inc", bson.D{{field, 1}}},
	}, options.Update().SetUpsert(true))

	if err == nil {return nil}

	log.Println(err)
	return err
}
