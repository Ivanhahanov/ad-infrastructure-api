package database

import (
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func GetMetricsNames() {
	promAddr := os.Getenv("PROMETHEUS")
	if promAddr == ""{
		promAddr = "http://localhost:9090"
	}
	urlAddr := promAddr + "/api/v1/query_range"
	req, err := http.NewRequest("GET", urlAddr, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	startTime, _ := GetStartTimeStamp()

	q := req.URL.Query()
	q.Add("query", "checker")
	q.Add("end", time.Now().Format(time.RFC3339))
	q.Add("start", startTime)
	q.Add("step", config.Conf.RoundInterval)
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())
	client := &http.Client{}
	resp, reqErr := client.Do(req)
	if reqErr != nil {
		log.Println(reqErr)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Println(string(body))
}
