package database

import (
	"github.com/go-redis/redis"
	"log"
	"time"
)

var client *redis.Client
var timeClient *redis.Client

func InitRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "admin",
		DB:       0,
	})
	timeClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "admin",
		DB:       1,
	})
}

type FlagStruct struct {
	Flag    string
	Team    string
	Service string
}

func PutFlag(flagStruct *FlagStruct) error {
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
func RemoveAllFlags() {
	client.FlushDB()
}

func WriteTime() {
	timeClient.RPush("time", time.Now().Format(time.RFC3339))
}
func GetTime(index int64) (string, error) {
	result, err := timeClient.LIndex("time", index).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}
func GetStartTimeStamp() (string, error) {
	return GetTime(0)
}

func GetLastTimeStamp() (string, error) {
	return GetTime(-1)
}

func WriteFlagSubmit(flagStruct *FlagStruct){
	status := client.HMSet(flagStruct.Flag, map[string]interface{}{
		"team":    flagStruct.Team,
		"service": flagStruct.Service,
	})
	log.Println(status)
}