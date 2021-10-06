package repositories

import "github.com/Ivanhahanov/ad-infrastructure-api/pkg/db"

type RepoManager struct {
	FlagRepo *FlagRepository
	TeamRepo *TeamRepository
}

func NewRepoManager(db *db.Db) *RepoManager {
	return &RepoManager{
		FlagRepo: NewFlagRepository(db),
		TeamRepo: NewTeamRepository(db),
	}
}

