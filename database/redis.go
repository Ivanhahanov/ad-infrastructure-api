package database

import (
	"github.com/go-redis/redis"
	"log"
)

var client *redis.Client

func InitRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "admin",
		DB:       0,
	})
}

type FlagStruct struct {
	Flag    string
	Team    string
	Service string
}

func PutFlag(flagStruct FlagStruct) error {
	status := client.HMSet(flagStruct.Flag, map[string]interface{}{
		"team":    flagStruct.Team,
		"service": flagStruct.Service,
	})
	log.Println(status)
	return nil
}

func GetInfo(flag string) ([]interface{}, error) {
	result, err := client.HMGet(flag, "team", "service").Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}
func RemoveAllFlags()  {
	client.FlushAll()
}
