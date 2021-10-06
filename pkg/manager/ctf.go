package manager

import (
	"github.com/Ivanhahanov/ad-infrastructure-api/config"
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/manager/repositories"
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/utils"
	"github.com/go-redis/redis"
	"time"
)

type CtfManager struct {
	FlagRepo *repositories.FlagRepository
	TeamRepo *repositories.TeamRepository

	client 		 *redis.Client
	timeClient   *redis.Client
	submitClient *redis.Client

	cfg *config.Config
}

func newRedisClient(addr string, pass string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})
}

func NewManager(repoManager *repositories.RepoManager, cfg *config.Config) *CtfManager {
	addr := utils.GetEnv("REDIS", "localhost:6379")
	pass := utils.GetEnv("ADMIN_PASS", "admin")

	return &CtfManager{
		FlagRepo: repoManager.FlagRepo,
		TeamRepo: repoManager.TeamRepo,

		client:       newRedisClient(addr, pass, 0),
		timeClient:   newRedisClient(addr, pass, 1),
		submitClient: newRedisClient(addr, pass, 2),

		cfg: cfg,
	}
}

func (m *CtfManager) PutFlag(flag string, team string, service string) (string, error)  {
	fields := map[string]interface{}{ "team": team, "service": service }
	result, err := m.client.HMSet(flag, fields).Result()
	if err != nil {return "", err}

	return result, nil
}

func (m *CtfManager) GetFlags(flag string) ([]interface{}, error) {
	result, err := m.client.HMGet(flag, "team", "service").Result()
	if err != nil {return nil, err}

	return result, nil
}

func (m *CtfManager) GetSubmitFlags(flag string) ([]interface{}, error) {
	result, err := m.submitClient.HMGet(flag, "team", "service").Result()
	if err != nil {return nil, err}

	return result, nil
}

func (m *CtfManager) GetTime(index int64) (string, error) {
	result, err := m.timeClient.LIndex("time", index).Result()
	if err != nil {return "", err}

	return result, nil
}

func (m *CtfManager) WriteTime() {
	m.timeClient.RPush("time", time.Now().Format(time.RFC3339))
}

func (m *CtfManager) GetStartTimeStamp() (string, error) {return m.GetTime(0)}
func (m *CtfManager) GetLastTimeStamp()  (string, error) {return m.GetTime(-1)}

func (m *CtfManager) RemoveAllFlags() {m.client.FlushDB()}