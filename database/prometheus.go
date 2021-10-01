package database

import (
	"encoding/json"
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type PrometheusQueryRange struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Name     string `json:"__name__"`
				Instance string `json:"instance"`
				Job      string `json:"job"`
				Proto    string `json:"proto"`
				Route    string `json:"route"`
				Service  string `json:"service"`
				Team     string `json:"team"`
			} `json:"metric"`
			Values [][]interface{} `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

func GetMetricsNames(teamName string, taskName string) {
	promAddr := os.Getenv("PROMETHEUS")
	if promAddr == "" {
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
	q.Add("query", fmt.Sprintf("{team=\"%s\",service=\"%s\"}",teamName, taskName ))
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
		log.Println(err)
		return
	}
	queryRanges := PrometheusQueryRange{}
	jsonErr := json.Unmarshal(body, &queryRanges)
	if jsonErr != nil {
		log.Println(jsonErr)
		return
	}
	sources := float64(len(queryRanges.Data.Result))
	if sources < 1{
		log.Println("No Data")
		return
	}
	rounds := float64(len(queryRanges.Data.Result[0].Values))
	sumValues := 0.0
	log.Println(queryRanges.Data.Result[0].Metric.Team, queryRanges.Data.Result[0].Metric.Service)
	for _, query := range queryRanges.Data.Result{
		for _, values := range query.Values{
			value, _ := strconv.Atoi(values[1].(string))
			sumValues += float64(value)
		}
	}
	log.Printf("%.2f%%", sumValues/rounds/sources*100)
}
