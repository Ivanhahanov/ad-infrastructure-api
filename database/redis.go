package database

import (
	"github.com/go-redis/redis"
	"log"
	"os"
	"time"
)

var client *redis.Client
var timeClient *redis.Client
var submitClient *redis.Client

func InitRedis() {
	redisAddr := os.Getenv("REDIS")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "admin",
		DB:       0,
	})
	submitClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "admin",
		DB:       1,
	})
	timeClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "admin",
		DB:       2,
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

func AddSubmitFlag(flagStruct *FlagStruct) {
	status := submitClient.HMSet(flagStruct.Flag, map[string]interface{}{
		"team":    flagStruct.Team,
		"service": flagStruct.Service,
	})
	log.Println(status)
}

func GetSubmitFlags(flag string) ([]interface{}, error) {
	result, err := submitClient.HMGet(flag, "team", "service").Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}
