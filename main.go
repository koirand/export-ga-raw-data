package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/api/analytics/v3"
	"google.golang.org/api/option"
)

const credentialsPath = "./credentials.json"
const viewID = "208412389"
const startDate = "90daysAgo"
const endDate = "yesterday"

func main() {
	s, err := getService()
	if err != nil {
		log.Fatalln(err)
	}

	results, err := s.Data.Ga.Get("ga:"+viewID, startDate, endDate, "ga:pageviews").
		Output("json").
		Dimensions(strings.Join([]string{
			"ga:clientId",
			"ga:pagePath",
			"ga:dateHourMinute",
		}, ",")).
		Do()
	if err != nil {
		log.Fatalln(err)
	}
	for _, row := range results.Rows {
		fmt.Println(strings.Join(row, ","))
	}
}

func getService() (*analytics.Service, error) {
	ctx := context.Background()
	as, err := analytics.NewService(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return nil, err
	}
	return as, nil
}
